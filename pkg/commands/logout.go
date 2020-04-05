package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/progress"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

// LogoutCmd find users in pool sub command
type LogoutCmd struct {
	UserPoolID string            `help:"User pool id." kong:"required"`
	BackOff    int               `help:"Delay in ms used to backoff during paging of records" default:"500"`
	Filter     map[string]string `help:"Filter users based on a set of patterns, supports  '*' and '?' wildcards in either string."`
}

// Run run the list operation
func (f *LogoutCmd) Run(ctx *Context) error {
	log.Debug().Msg("logout users")

	log.Debug().Fields(convertMap(f.Filter)).Msg("Filter")

	filteringEnabled := len(f.Filter) > 0

	// use this to keep a track of how many entries we will update and logout
	usernames := []string{}

	err := ctx.Cognito.ListUsers(f.UserPoolID, func(p *cognito.UsersPage) bool {
		log.Debug().Int("len", len(p.Users)).Msg("page")

		for _, user := range p.Users {

			m := attrToMap(user.Attributes)

			m["Username"] = aws.StringValue(user.Username)

			if filteringEnabled && !matchFilters(f.Filter, m) {
				continue
			}

			usernames = append(usernames, aws.StringValue(user.Username))

		}

		time.Sleep(time.Duration(f.BackOff) * time.Millisecond)

		return true // continue paging
	})
	if err != nil {
		return err
	}

	log.Debug().Int("len", len(usernames)).Msg("users will be logged out")

	if len(usernames) == 0 {
		fmt.Fprintln(ctx.Writer, "No users found.")
		return nil
	}

	fmt.Fprintf(ctx.Writer, "Found users commencing logout for count=%d\n", len(usernames))

	time.Sleep(1 * time.Second)

	pw := progress.NewWriter()
	pw.SetOutputWriter(os.Stderr)

	tracker := &progress.Tracker{Message: "users being logged out", Total: int64(len(usernames))}

	pw.AppendTracker(tracker)

	go pw.Render()

	for _, username := range usernames {
		log.Debug().Str("username", username).Msg("calling logout for user")

		err := ctx.Cognito.Logout(f.UserPoolID, username)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to logout user") // best to just stop here
		}

		tracker.Increment(1)
		time.Sleep(time.Duration(f.BackOff) * time.Millisecond)
	}

	return nil
}
