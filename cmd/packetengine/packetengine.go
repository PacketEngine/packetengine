package main

import (
	"errors"
	"os"

	"github.com/PacketEngine/packetengine"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var colorError *color.Color
var colorInfo *color.Color

var packetengineClient *packetengine.PacketEngineClient

var withoutTags string
var allSubdomains bool

func main() {
	colorError = color.New(color.FgWhite).Add(color.BgRed)
	colorInfo = color.New(color.FgGreen)

	app := &cli.App{
		Name:    "PacketEngine CLI",
		Usage:   "CLI application for the PacketEngine (packetengine.co.uk) API",
		Suggest: true,
		Commands: []*cli.Command{
			{
				Name:        "init",
				HelpName:    "init",
				Action:      InitAction,
				ArgsUsage:   ` `,
				Usage:       `Sets the application API token.`,
				Description: `Set the API token.`,
			},
			{
				Name:        "subdomains",
				HelpName:    "subdomains",
				Action:      SubdomainsAction,
				ArgsUsage:   ` `,
				Usage:       `Fetches subdomains of a domain.`,
				Description: `Fetches subdomains of a domain.`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "without-tags",
						Usage:       "Without tags",
						Destination: &withoutTags,
					},
					&cli.BoolFlag{
						Name:        "all",
						Usage:       "If set to true, include subdomains which don't have any current DNS records",
						Destination: &allSubdomains,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		colorError.Println(err)
	}
}

func InitAction(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("An API token is required.")
	}

	token := c.Args().First()

	err := SetAPIToken(token)
	if err != nil {
		return err
	}

	colorInfo.Println("API token set!")
	return nil
}

func SubdomainsAction(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("A domain is required.")
	}

	// Boot
	apiToken, err := GetAPIToken()

	if err != nil {
		return err
	}

	packetengineClient, err = packetengine.NewPacketEngineClient(apiToken)

	if err != nil {
		return err
	}

	domain := c.Args().First()

	subdomains, err := packetengineClient.GetSubdomains(domain, withoutTags, allSubdomains)
	if err != nil {
		return err
	}

	for _, subdomain := range subdomains {
		colorInfo.Println(subdomain)
	}

	return nil
}

func SetAPIToken(apiToken string) error {
	if len(apiToken) == 0 {
		return errors.New("An API token is required.")
	}

	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	tokenBytes := []byte(apiToken)
	err = os.MkdirAll(home+"/.config/packetengine/", os.ModePerm)

	if err != nil {
		return err
	}

	err = os.WriteFile(home+"/.config/packetengine/api_key", tokenBytes, 0644)

	if err != nil {
		return err
	}

	return nil
}

func GetAPIToken() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	apiToken, err := os.ReadFile(home + "/.config/packetengine/api_key")

	if err != nil {
		return "", err
	}

	return string(apiToken), nil
}
