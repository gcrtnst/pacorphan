package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type MakePkg struct {
	BinPath string
	CfgPath string
	Stdout  io.Writer
	Stderr  io.Writer
}

func NewMakePkg() (*MakePkg, error) {
	var err error

	bin, err := exec.LookPath("makepkg")
	if err != nil {
		return nil, err
	}

	cfg, err := CreateTempMakePkgConf()
	if err != nil {
		return nil, err
	}

	m := &MakePkg{
		BinPath: bin,
		CfgPath: cfg,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
	return m, nil
}

func (m *MakePkg) Dispose() error {
	return os.Remove(m.CfgPath)
}

func (m *MakePkg) Run(dst string, src *PkgBuild) ([]string, error) {
	bin := m.BinPath
	var err error

	dst, err = filepath.Abs(dst)
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // ignore error

	srcName := filepath.Join(dir, "PKGBUILD")
	err = src.WriteFile(srcName)
	if err != nil {
		return nil, err
	}

	arg := []string{
		bin, "--config", m.CfgPath,
		"-p", srcName, "--nodeps", "--clean", "--cleanbuild",
	}

	var env []string
	for _, e := range os.Environ() {
		if name, _, found := strings.Cut(e, "="); found {
			switch name {
			case "CHOST", "CARCH", "PACKAGER", "BUILDENV", "OPTIONS", "PKGEXT", "SRCEXT":
				continue
			}
		}
		env = append(env, e)
	}
	env = append(
		env,
		"PKGDEST="+dst,
		"SRCDEST="+dir,
		"SRCPKGDEST="+dir,
		"LOGDEST="+dir,
		"BUILDDIR="+dir,
	)

	cmdList := &exec.Cmd{
		Path: bin,
		Args: append(slices.Clone(arg), "--packagelist"),
		Env:  env,
		Dir:  dir,
	}
	outList, err := cmdList.Output()
	if err != nil {
		return nil, WrapExitError(cmdList, err)
	}

	pkgList := make([]string, 0, 1)
	for out := range bytes.SplitSeq(outList, []byte("\n")) {
		if len(out) > 0 {
			pkg := string(out)
			if !filepath.IsAbs(pkg) {
				pkg = filepath.Join(dir, pkg)
			}
			pkgList = append(pkgList, pkg)
		}
	}

	cmdMake := &exec.Cmd{
		Path:   bin,
		Args:   arg,
		Env:    env,
		Dir:    dir,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
	err = cmdMake.Run()
	if err != nil {
		return nil, WrapExitError(cmdMake, err)
	}

	return pkgList, nil
}

func CreateTempMakePkgConf() (_ string, err error) {
	arch, err := Arch()
	if err != nil {
		return "", err
	}

	b := BashScript{
		&BashVar{Name: "CARCH", Value: arch},
		&BashVar{Name: "PKGEXT", Value: ".pkg.tar.gz"},
		&BashVar{Name: "SRCEXT", Value: ".src.tar.gz"},
	}

	f, err := os.CreateTemp("", "makepkg-*.conf")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
		if err != nil {
			_ = os.Remove(f.Name())
		}
	}()

	_, err = f.Write([]byte(b.Text()))
	if err != nil {
		return "", err
	}

	err = f.Sync()
	if err != nil {
		return "", err
	}

	name := f.Name()
	name, err = filepath.Abs(name)
	if err != nil {
		return "", err
	}

	return name, nil
}

type PkgBuild struct {
	Name       string
	Ver        string
	Rel        int
	Arch       []string
	Depends    []string
	OptDepends []string
}

func NewPkgBuild(name string, ver string) *PkgBuild {
	return &PkgBuild{
		Name: name,
		Ver:  ver,
		Rel:  1,
		Arch: []string{"any"},
	}
}

func (p *PkgBuild) WriteFile(name string) error {
	b := BashScript{
		&BashVar{Name: "pkgname", Value: p.Name},
		&BashVar{Name: "pkgver", Value: p.Ver},
		&BashInt{Name: "pkgrel", Value: p.Rel},
		&BashArray{Name: "arch", Value: p.Arch},
		&BashArray{Name: "depends", Value: p.Depends, Optional: true},
		&BashArray{Name: "optdepends", Value: p.OptDepends, Optional: true},
	}
	return b.WriteFile(name)
}

type BashScript []BashFragment

func (b BashScript) WriteFile(name string) error {
	s := b.Text()
	return os.WriteFile(name, []byte(s), 0o644)
}

func (b BashScript) Text() string {
	s := &strings.Builder{}
	for _, f := range b {
		f.AppendBashFragment(s)
	}
	return s.String()
}

type BashFragment interface {
	AppendBashFragment(*strings.Builder)
}

type BashVar struct {
	Name  string
	Value string
}

func (b *BashVar) AppendBashFragment(s *strings.Builder) {
	// validation of name is omitted

	_, _ = s.WriteString(b.Name)
	_, _ = s.WriteRune('=')
	_, _ = s.WriteString(QuoteBashString(b.Value))
	_, _ = s.WriteRune('\n')
}

type BashInt struct {
	Name  string
	Value int
}

func (b *BashInt) AppendBashFragment(s *strings.Builder) {
	// validation of name is omitted

	_, _ = s.WriteString(b.Name)
	_, _ = s.WriteRune('=')
	_, _ = s.WriteString(strconv.Itoa(b.Value))
	_, _ = s.WriteRune('\n')
}

type BashArray struct {
	Name     string
	Value    []string
	Optional bool
}

func (b *BashArray) AppendBashFragment(s *strings.Builder) {
	if b.Optional && len(b.Value) <= 0 {
		return
	}

	// validation of name is omitted

	qValue := make([]string, 0, len(b.Value))
	for _, elem := range b.Value {
		qElem := QuoteBashString(elem)
		qValue = append(qValue, qElem)
	}

	_, _ = s.WriteString(b.Name)
	_, _ = s.WriteString("=(")
	_, _ = s.WriteString(strings.Join(qValue, " "))
	_, _ = s.WriteString(")\n")
}

func Arch() (string, error) {
	cmd := exec.Command("uname", "--machine")
	out, err := cmd.Output()
	if err != nil {
		return "", WrapExitError(cmd, err)
	}

	arch := strings.TrimSpace(string(out))
	return arch, nil
}

func QuoteBashString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
