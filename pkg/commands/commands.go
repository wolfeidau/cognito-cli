package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jedib0t/go-pretty/table"
	"github.com/rs/zerolog/log"
)

// Context cli context used for common options
type Context struct {
	Debug bool
}

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

// FindCmd find users in pool sub command
type FindCmd struct {
	UserPoolID string   `help:"User pool id." kong:"required"`
	CSV        bool     `help:"Enable csv output."`
	Attributes []string `help:"Attributes to retrieve and output." default:"email"`
	BackOff    int      `help:"Delay in ms used to backoff during paging of records" default:"500"`
}

// Run run the list operation
func (f *FindCmd) Run(ctx *Context) error {
	log.Debug().Strs("attr", f.Attributes).Msg("find users")

	sess := session.Must(session.NewSession())
	csvc := cognitoidentityprovider.New(sess)

	tw := table.NewWriter()

	tw.AppendHeader(buildTableHeader(f.Attributes))

	err := csvc.ListUsersPagesWithContext(context.TODO(), &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(f.UserPoolID),
		Limit:      aws.Int64(10),
	},
		func(p *cognitoidentityprovider.ListUsersOutput, lastPage bool) bool {
			log.Debug().Int("len", len(p.Users)).Msg("page")

			for _, user := range p.Users {

				m := attrToMap(user.Attributes)

				tr := table.Row{}

				tr = append(tr, aws.StringValue(user.Username))

				for _, att := range f.Attributes {
					tr = append(tr, m[att])

				}

				tr = append(tr, aws.BoolValue(user.Enabled))
				tr = append(tr, aws.TimeValue(user.UserLastModifiedDate).Local())

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
		fmt.Println("No users found.")
		return nil
	}

	tw.SortBy([]table.SortBy{{Name: "LastModified", Mode: table.Dsc}})

	if f.CSV {
		fmt.Println(tw.RenderCSV())
	} else {
		fmt.Println(tw.Render())
	}

	return nil
}

func attrToMap(attributes []*cognitoidentityprovider.AttributeType) map[string]string {
	m := map[string]string{}

	for _, attr := range attributes {
		m[aws.StringValue(attr.Name)] = aws.StringValue(attr.Value)
	}

	return m
}

func buildTableHeader(names []string) table.Row {
	h := table.Row{}
	h = append(h, "Username")

	for _, name := range names {
		h = append(h, name)
	}

	h = append(h, "Enabled")
	h = append(h, "LastModified")

	return h
}
