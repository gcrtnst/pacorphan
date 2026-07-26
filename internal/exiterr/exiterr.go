package exiterr

import (
	"errors"
	"os/exec"
	"strings"
)

type ExitError struct {
	Cmd *exec.Cmd
	Err *exec.ExitError
}

func Wrap(cmd *exec.Cmd, err error) error {
	if errExit, ok := errors.AsType[*exec.ExitError](err); ok {
		return &ExitError{Cmd: cmd, Err: errExit}
	}
	return err
}

func (e *ExitError) Error() string {
	if e == nil {
		return "<nil>"
	}

	msgCmd := "exec"
	if e.Cmd != nil {
		msgCmd = "exec [" + strings.Join(e.Cmd.Args, " ") + "]"
	}

	msgErr := "<nil>"
	if e.Err != nil {
		msgErr = e.Err.Error()
	}

	return msgCmd + ": " + msgErr
}

func (e *ExitError) Unwrap() error {
	if e == nil || e.Err == nil {
		return nil
	}
	return e.Err
}
