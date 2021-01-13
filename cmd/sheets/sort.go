package sheets

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type SortOptions struct {
	*CommandOptions
	SpreadSheetID string
}

// NewCmdSort represents the sort command
func NewCmdSort(c *CommandOptions) *cobra.Command {
	opts := &SortOptions{CommandOptions: c}
	var cmd = &cobra.Command{
		Use:   "sort <spread-sheet-id>",
		Short: "Sort a sheets",
		Example: heredoc.Doc(`
		$ gospread sheets sort abc1234
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.SpreadSheetID = args[0]
			}
			return sortRun(opts)
		},
	}
	return cmd
}

func sortRun(opts *SortOptions) (err error) {
	g := opts.GoSpread
	spreadID := spreadIDFromArg(opts.SpreadSheetID)
	if err := g.SortSheets(spreadID); err != nil {
		return err
	}
	return nil
}
