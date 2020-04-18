package commands

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
)

// ListAttributesCmd list pools attributes sub command
type ListAttributesCmd struct {
	UserPoolID string `help:"User pool id." kong:"required"`
}

// Run run the list operation
func (f *ListAttributesCmd) Run(cli *CLIContext) error {

	log.Debug().Msg("list attributes")

	attrs, err := cli.Cognito.DescribePoolAttributes(f.UserPoolID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list pool attributes")
	}

	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"name"})

	for _, att := range attrs {
		tw.AppendRow(table.Row{att})
	}

	fmt.Fprintln(cli.Writer, tw.Render())

	return nil
}
