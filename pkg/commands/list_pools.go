package commands

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
)

// LsCmd list pools sub command
type LsCmd struct {
	CSV bool `help:"Enable csv output."`
}

// Run run the list operation
func (l *LsCmd) Run(ctx *Context) error {
	log.Debug().Msg("list pools")

	sess := session.Must(session.NewSession())
	csvc := cognitoidentityprovider.New(sess)

	tw := table.NewWriter()

	tw.AppendHeader(table.Row{"ID", "Name", "Created"})

	err := csvc.ListUserPoolsPagesWithContext(context.TODO(), &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(60),
	},
		func(p *cognitoidentityprovider.ListUserPoolsOutput, lastPage bool) bool {
			log.Debug().Int("len", len(p.UserPools)).Msg("page")

			for _, pool := range p.UserPools {
				tw.AppendRows([]table.Row{{aws.StringValue(pool.Id), aws.StringValue(pool.Name), aws.TimeValue(pool.CreationDate).Local()}})
			}

			return true // continue paging
		})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list pools")
	}

	log.Debug().Int("len", tw.Length()).Msg("render table")

	if tw.Length() == 0 {
		fmt.Println("No pools found.")
		return nil
	}

	tw.SortBy([]table.SortBy{{Name: "Created", Mode: table.Dsc}})

	if l.CSV {
		fmt.Println(tw.RenderCSV())
	} else {
		fmt.Println(tw.Render())
	}

	return nil
}
