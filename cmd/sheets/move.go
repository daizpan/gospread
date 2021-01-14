package sheets

import (
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type MoveOptions struct {
	*CommandOptions
	SpreadSheetID string
	SheetID       int64
	Index         int64
}

// NewCmdMove represents the move command
func NewCmdMove(c *CommandOptions) *cobra.Command {
	opts := &MoveOptions{CommandOptions: c}
	var cmd = &cobra.Command{
		Use:   "move <sheet-id> <index> {<spread-id> | <url>}",
		Short: "Move a sheet",
		Example: heredoc.Doc(`
		$ gospread sheets move sheet1 0 abc1234
		`),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				id, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
				opts.SheetID = id
				index, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return err
				}
				opts.Index = index
				opts.SpreadSheetID = args[2]
			}
			return moveRun(opts)
		},
	}
	return cmd
}

func moveRun(opts *MoveOptions) (err error) {
	g := opts.GoSpread
	spreadID := spreadIDFromArg(opts.SpreadSheetID)
	if err := g.MoveSheet(opts.SheetID, opts.Index, spreadID); err != nil {
		return err
	}
	return nil
}
