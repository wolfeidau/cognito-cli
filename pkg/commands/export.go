package commands

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

// ExportCmd find and export users in pool sub command
type ExportCmd struct {
	UserPoolID string            `help:"User pool id." kong:"required"`
	BackOff    int               `help:"Delay in ms used to backoff during paging of records" default:"500"`
	Filter     map[string]string `help:"Filter users based on a set of patterns, supports  '*' and '?' wildcards in either string."`
}

// Run run the list operation
func (f *ExportCmd) Run(ctx *Context) error {
	log.Debug().Msg("find and export users")

	tw := table.NewWriter()

	attrs, err := ctx.Cognito.DescribePoolAttributes(f.UserPoolID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list pool attributes")
	}

	// prepend the Username
	attrs = append([]string{"Username"}, attrs...)

	tw.AppendHeader(buildTableHeader(attrs))

	log.Debug().Fields(convertMap(f.Filter)).Msg("Filter")

	filteringEnabled := len(f.Filter) > 0

	err = ctx.Cognito.ListUsers(f.UserPoolID, func(p *cognito.UsersPage) bool {
		log.Debug().Int("len", len(p.Users)).Msg("page")

		for _, user := range p.Users {

			m := attrToMap(user.Attributes)
			m["Username"] = aws.StringValue(user.Username)

			if filteringEnabled && !matchFilters(f.Filter, m) {
				continue
			}

			tr := table.Row{}

			for _, attr := range attrs {
				if _, ok := m[attr]; ok {
					tr = append(tr, m[attr])
				} else {
					tr = append(tr, "")
				}
			}

			tr = append(tr, aws.BoolValue(user.Enabled))
			tr = append(tr, awsTimeLocal(user.UserLastModifiedDate, !ctx.DisableLocalTime))

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
		fmt.Fprintln(ctx.Writer, "No users found.")
		return nil
	}

	tw.SortBy([]table.SortBy{{Name: "LastModified", Mode: table.Dsc}})

	fmt.Fprintln(ctx.Writer, tw.RenderCSV())

	return nil
}
