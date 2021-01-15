package sheets

import (
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	*CommandOptions
	SpreadSheetID string
	Exclusions    []string
}

// NewCmdDeleteall represents the deleteall command
func NewCmdDeleteall(c *CommandOptions) *cobra.Command {
	opts := &DeleteOptions{CommandOptions: c}
	var cmd = &cobra.Command{
		Use:   "deleteall <spread-id>",
		Short: "Delete all sheets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.SpreadSheetID = args[0]
			}
			return deleteAllRun(opts)
		},
	}
	fl := cmd.Flags()
	fl.StringSliceVarP(&opts.Exclusions, "exclusion", "e", nil, "Exclusions sheet name.")
	return cmd
}

func deleteAllRun(opts *DeleteOptions) error {
	g := opts.GoSpread
	spreadID := spreadIDFromArg(opts.SpreadSheetID)
	if err := g.DeleteAllSheet(spreadID, opts.Exclusions); err != nil {
		return err
	}
	return nil
}
