package base

import (
	"fmt"
	"github.com/go-cmd/cmd"
)

type Cmd struct {
	ignoreStdoutStdErr bool
	waitMode           bool
	successCode        int
}

func NewCmd() *Cmd {
	return &Cmd{
		ignoreStdoutStdErr: true,
		waitMode:           false,
		successCode:        0,
	}
}

func (c *Cmd) SetSuccessCode(code int) {
	c.successCode = code
}

func (c *Cmd) Run(name string, args []string) error {
	cmd := cmd.NewCmd(name, args...)
	status := <-cmd.Start()
	if status.StartTs > 0 && (status.Exit != c.successCode || !status.Complete) {
		return fmt.Errorf("%d", status.Exit)
	}
	if status.StartTs == 0 && status.Error != nil {
		return status.Error
	}
	return nil
}
