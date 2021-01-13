package cmd

import (
	"os"

	"github.com/daizpan/gospread"
	"github.com/daizpan/gospread/cmd/sheets"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
func NewCmdRoot() *cobra.Command {
	opts := &sheets.CommandOptions{}
	cmd := &cobra.Command{
		Version: "0.1.0",
		Use:     "gospread",
		Short:   "Manage a Google Spread Sheet",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			g, err := gospread.NewGoSpreadWithCredentialFile(opts.Credential)
			if err != nil {
				return err
			}
			opts.GoSpread = g
			return nil
		},
	}
	cmd.AddCommand(
		sheets.NewCmdSheets(opts),
	)
	fl := cmd.PersistentFlags()
	fl.StringVarP(&opts.Credential, "credential", "c", "credential.json", "Sheets api credential json file")
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewCmdRoot()
	cmd.SetOutput(os.Stdout)
	cmd.SetErr(os.Stderr)

	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
