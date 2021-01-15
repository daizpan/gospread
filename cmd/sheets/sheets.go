package sheets

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// NewCmdSheets represents the sheets command
func NewCmdSheets(c *CommandOptions) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sheets <command>",
		Short: "Manage spread sheets",
		Example: heredoc.Doc(`
		$ gospread sheets create <sheet-name> <spread-id>
		$ gospread sheets sort <spread-id>
		$ gospread sheets move <sheet-nameid> <index> <spread-id>
		`),
	}

	cmd.AddCommand(
		NewCmdCreate(c),
		NewCmdSort(c),
		NewCmdMove(c),
		NewCmdDeleteall(c),
	)
	return cmd
}
