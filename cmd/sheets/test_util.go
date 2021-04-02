package sheets

import (
	"bytes"
	"os"

	"github.com/duhshu/gospread"
	"github.com/spf13/cobra"
)

type CmdOut struct {
	OutBuf, ErrBuf *bytes.Buffer
}

func (c CmdOut) String() string {
	return c.OutBuf.String()
}

func (c CmdOut) Stderr() string {
	return c.ErrBuf.String()
}

func runCommand(newCmdFunc func(c *CommandOptions) *cobra.Command, args ...string) (*CmdOut, error) {

	g, err := gospread.NewGoSpreadWithCredentialFile(os.Getenv("SHEET_API_CREDENTIAL_JSON"))
	if err != nil {
		return nil, err
	}
	opts := &CommandOptions{GoSpread: g}
	cmd := newCmdFunc(opts)

	// argv, err := shlex.Split(cli)
	// if err != nil {
	// 	return nil, err
	// }
	cmd.SetArgs(args)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.SetOutput(outBuf)
	cmd.SetErr(errBuf)
	err = cmd.Execute()
	return &CmdOut{
		OutBuf: outBuf,
		ErrBuf: errBuf,
	}, err
}
