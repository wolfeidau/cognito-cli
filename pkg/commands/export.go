package commands

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

	csvw := csv.NewWriter(ctx.Writer)

	attrs, err := ctx.Cognito.DescribePoolAttributes(f.UserPoolID)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list pool attributes")
	}

	// prepend the Username
	attrs = append([]string{"Username"}, attrs...)

	err = csvw.Write(append(attrs, "Enabled", "LastModified"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to write headers for CSV")
	}

	log.Debug().Fields(convertMap(f.Filter)).Msg("Filter")

	filteringEnabled := len(f.Filter) > 0

	count := 0

	err = ctx.Cognito.ListUsers(f.UserPoolID, func(p *cognito.UsersPage) bool {
		log.Debug().Int("len", len(p.Users)).Msg("page")

		for _, user := range p.Users {

			m := attrToMap(user.Attributes)
			m["Username"] = aws.StringValue(user.Username)

			if filteringEnabled && !matchFilters(f.Filter, m) {
				continue
			}

			tr := []string{}

			for _, attr := range attrs {
				if _, ok := m[attr]; ok {
					tr = append(tr, m[attr])
				} else {
					tr = append(tr, "")
				}
			}

			tr = append(tr, fmt.Sprintf("%t", aws.BoolValue(user.Enabled)))
			tr = append(tr, awsTimeLocal(user.UserLastModifiedDate, !ctx.DisableLocalTime).String())

			err = csvw.Write(tr)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to write row to CSV")
			}
			count++
		}

		time.Sleep(time.Duration(f.BackOff) * time.Millisecond)

		return true // continue paging
	})
	if err != nil {
		return err
	}

	log.Debug().Int("len", count).Msg("render table")

	if count == 0 {
		fmt.Fprintln(ctx.Writer, "No users found.")
		return nil
	}

	// Write any buffered data to the underlying writer (standard output).
	csvw.Flush()

	if err = csvw.Error(); err != nil {
		log.Fatal().Err(err).Msg("failed to write flush data to CSV")
	}

	return nil
}
