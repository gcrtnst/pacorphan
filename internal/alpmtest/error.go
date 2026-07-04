package main

import (
	"errors"
	"os/exec"
	"strings"
)

type ExitError struct {
	Cmd *exec.Cmd
	Err *exec.ExitError
}

func WrapExitError(cmd *exec.Cmd, err error) error {
	if errExit, ok := errors.AsType[*exec.ExitError](err); ok {
		return &ExitError{Cmd: cmd, Err: errExit}
	}
	return err
}

func (e *ExitError) Error() string {
	msgCmd := "exec <nil>"
	if e.Cmd != nil {
		msgCmd = "exec [" + strings.Join(e.Cmd.Args, " ") + "]"
	}

	msgErr := "unknown error"
	if e.Err != nil {
		msgErr = e.Err.Error()
	}

	return msgCmd + ": " + msgErr
}

func (e *ExitError) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return e.Err
}
