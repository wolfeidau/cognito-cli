package commands

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

// FindCmd find users in pool sub command
type FindCmd struct {
	UserPoolID string            `help:"User pool id." kong:"required"`
	CSV        bool              `help:"Enable csv output."`
	Attributes []string          `help:"Attributes to retrieve and output." default:"Username,email"`
	BackOff    int               `help:"Delay in ms used to backoff during paging of records" default:"500"`
	Filter     map[string]string `help:"Filter users based on a set of patterns, supports  '*' and '?' wildcards in either string."`
}

// Run run the list operation
func (f *FindCmd) Run(cli *CLIContext) error {
	log.Debug().Strs("attr", f.Attributes).Msg("find users")

	tw := table.NewWriter()

	tw.AppendHeader(buildTableHeader(f.Attributes))

	log.Debug().Fields(convertMap(f.Filter)).Msg("Filter")

	filteringEnabled := len(f.Filter) > 0

	err := cli.Cognito.ListUsers(f.UserPoolID, func(p *cognito.UsersPage) bool {
		log.Debug().Int("len", len(p.Users)).Msg("page")

		for _, user := range p.Users {

			m := attrToMap(user.Attributes)
			m["Username"] = aws.StringValue(user.Username)

			if filteringEnabled && !matchFilters(f.Filter, m) {
				continue
			}

			tr := table.Row{}

			for _, att := range f.Attributes {
				tr = append(tr, m[att])

			}

			tr = append(tr, aws.BoolValue(user.Enabled))
			tr = append(tr, awsTimeLocal(user.UserLastModifiedDate, !cli.DisableLocalTime))

			tw.AppendRows([]table.Row{tr})
		}

		time.Sleep(time.Duration(f.BackOff) * time.Millisecond)

		return true // continue paging
	})
	if err != nil {
		return err
	}

	log.Debug().Int("len", tw.Length()).Msg("render table")

	if tw.Length() == 0 {
		fmt.Fprintln(cli.Writer, "No users found.")
		return nil
	}

	tw.SortBy([]table.SortBy{{Name: "LastModified", Mode: table.Dsc}})

	if f.CSV {
		fmt.Fprintln(cli.Writer, tw.RenderCSV())
	} else {
		fmt.Fprintln(cli.Writer, tw.Render())
	}

	return nil
}
