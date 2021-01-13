package sheets

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	*CommandOptions
	Title         string
	SpreadSheetID string
}

// NewCmdCreate represents the create command
func NewCmdCreate(c *CommandOptions) *cobra.Command {
	opts := &CreateOptions{CommandOptions: c}
	var cmd = &cobra.Command{
		Use:   "create <sheet-name> <spread-id>",
		Short: "Create a sheet",
		Example: heredoc.Doc(`
		$ gospread sheets create NewSheet abc1234
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Title = args[0]
				opts.SpreadSheetID = args[1]
			}
			return createRun(opts)
		},
	}
	return cmd
}

func createRun(opts *CreateOptions) error {
	g := opts.GoSpread
	spreadID := spreadIDFromArg(opts.SpreadSheetID)
	if err := g.CreateSheet(spreadID, opts.Title); err != nil {
		return err
	}
	return nil
}
