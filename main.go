package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/internal/commands"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

var (
	version = "unknown"
)

var flags struct {
	Debug            bool `help:"Enable debug mode."`
	DisableLocalTime bool `help:"Disable localisation of times output."`
	Version          kong.VersionFlag

	Ls             commands.LsCmd             `cmd:"ls" help:"List pools."`
	ListAttributes commands.ListAttributesCmd `cmd:"list-attributes" help:"List the schema attributes of the user pool."`
	Find           commands.FindCmd           `cmd:"find" help:"Find users."`
	Export         commands.ExportCmd         `cmd:"export" help:"Export users, filter and write the results in CSV format."`
	Logout         commands.LogoutCmd         `cmd:"logout" help:"Find users and trigger a logout."`
}

func main() {
	cli := kong.Parse(&flags,
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

	if flags.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	cognitoSvc := cognito.NewService()

	err := cli.Run(&commands.CLIContext{Debug: flags.Debug, DisableLocalTime: flags.DisableLocalTime, Cognito: cognitoSvc, Writer: os.Stdout})
	cli.FatalIfErrorf(err)
}
