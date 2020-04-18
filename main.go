package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/cognito-cli/internal/commands"
	"github.com/wolfeidau/cognito-cli/pkg/cognito"
)

var (
	version = "unknown"
)

type regionFlag string

func (r regionFlag) AfterApply(cfg *aws.Config) error {
	cfg.Region = aws.String(string(r))
	return nil
}

type profileFlag string

func (p profileFlag) AfterApply(cfg *aws.Config) error {
	cfg.Credentials = credentials.NewSharedCredentials("", string(p))
	return nil
}

var flags struct {
	Debug            bool        `help:"Enable debug mode."`
	Region           regionFlag  `help:"AWS Region." env:"AWS_REGION" default:"us-east-1" short:"r"`
	Profile          profileFlag `help:"AWS CLI profile." env:"AWS_PROFILE" short:"p"`
	DisableLocalTime bool        `help:"Disable localisation of times output."`
	Version          kong.VersionFlag

	Ls             commands.LsCmd             `cmd:"ls" help:"List pools."`
	ListAttributes commands.ListAttributesCmd `cmd:"list-attributes" help:"List the schema attributes of the user pool."`
	Find           commands.FindCmd           `cmd:"find" help:"Find users."`
	Export         commands.ExportCmd         `cmd:"export" help:"Export users, filter and write the results in CSV format."`
	Logout         commands.LogoutCmd         `cmd:"logout" help:"Find users and trigger a logout."`
}

func main() {

	awscfg := &aws.Config{}

	cli := kong.Parse(&flags,
		kong.Bind(awscfg),
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

	log.Debug().Str("region", string(flags.Region)).Str("profile", string(flags.Profile)).Msg("aws config")

	cognitoSvc := cognito.NewService(awscfg)

	err := cli.Run(&commands.CLIContext{Debug: flags.Debug, DisableLocalTime: flags.DisableLocalTime, Cognito: cognitoSvc, Writer: os.Stdout})
	cli.FatalIfErrorf(err)
}
