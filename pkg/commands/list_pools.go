package commands

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

// LsCmd list pools sub command
type LsCmd struct {
	CSV bool `help:"Enable csv output."`
}

// Run run the list operation
func (l *LsCmd) Run(ctx *Context) error {
	log.Debug().Msg("list pools")

	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"ID", "Name", "Created"})

	err := ctx.Cognito.ListPools(func(p *cognito.UserPoolsPage) bool {
		log.Debug().Int("len", len(p.UserPools)).Msg("page")

		for _, pool := range p.UserPools {
			tw.AppendRows([]table.Row{{aws.StringValue(pool.Id), aws.StringValue(pool.Name), awsTimeLocal(pool.CreationDate, !ctx.DisableLocalTime)}})
		}

		return true // continue paging
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list pools")
	}

	log.Debug().Int("len", tw.Length()).Msg("render table")

	if tw.Length() == 0 {
		fmt.Fprintln(ctx.Writer, "No pools found.")
		return nil
	}

	tw.SortBy([]table.SortBy{{Name: "Created", Mode: table.Dsc}})

	if l.CSV {
		fmt.Fprintln(ctx.Writer, tw.RenderCSV())
	} else {
		fmt.Fprintln(ctx.Writer, tw.Render())
	}

	return nil
}
