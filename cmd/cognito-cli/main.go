package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
	"github.com/wolfeidau/cognito-cli/pkg/commands"
)

var (
	version = "unknown"
)

var cli struct {
	Debug   bool `help:"Enable debug mode."`
	Version kong.VersionFlag

	Ls   commands.LsCmd   `cmd:"ls" help:"List pools."`
	Find commands.FindCmd `cmd:"find" help:"Find users."`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Vars{"version": version}, // bind a var for version
		kong.Name("cognito-cli"),
		kong.Description("A cognito cli."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
			Summary: true,
		}))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cli.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cognitoSvc := cognito.NewService()

	err := ctx.Run(&commands.Context{Debug: cli.Debug, Cognito: cognitoSvc})
	ctx.FatalIfErrorf(err)
}
