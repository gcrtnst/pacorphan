package main

import (
	"bytes"

	"github.com/gcrtnst/pacorphan/internal/testenv"
)

func HelpEnv(t *T) *testenv.Env {
	return HelpEnvWithOption(t, nil)
}

func HelpEnvWithOption(t *T, opt *testenv.EnvOption) *testenv.Env {
	env, err := testenv.NewEnv(opt)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := env.Dispose()
		if err != nil {
			t.Log(err)
		}
	})

	env.SetStderr(t.Output())
	return env
}

func HelpMakeAndInstall(t *T, env *testenv.Env, src *testenv.PkgBuild, explicit bool) {
	we1 := env.Stdout
	we2 := env.Stderr
	wm1 := env.MakePkg.Stdout
	wm2 := env.MakePkg.Stderr
	defer func() {
		env.Stdout = we1
		env.Stderr = we2
		env.MakePkg.Stdout = wm1
		env.MakePkg.Stderr = wm2
	}()

	buf := new(bytes.Buffer)
	env.SetStdout(nil)
	env.SetStderr(buf)

	err := env.MakeAndInstall(src, explicit)
	if err != nil {
		t.Output().Write(buf.Bytes())
		t.Fatal(err)
	}
}
