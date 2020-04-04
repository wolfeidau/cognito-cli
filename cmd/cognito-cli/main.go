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
	Debug            bool `help:"Enable debug mode."`
	DisableLocalTime bool `help:"Disable localisation of times output."`
	Version          kong.VersionFlag

	Ls     commands.LsCmd     `cmd:"ls" help:"List pools."`
	Find   commands.FindCmd   `cmd:"find" help:"Find users."`
	Export commands.ExportCmd `cmd:"export" help:"Find users and export in CSV format."`
	Logout commands.LogoutCmd `cmd:"logout" help:"Find users and trigger a logout."`
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

	err := ctx.Run(&commands.Context{Debug: cli.Debug, DisableLocalTime: cli.DisableLocalTime, Cognito: cognitoSvc, Writer: os.Stdout})
	ctx.FatalIfErrorf(err)
}
