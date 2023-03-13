package main

import (
    "errors"
    "os"

    "github.com/fatih/color"
    "github.com/urfave/cli/v2"
    "github.com/PacketEngine/packetengine"
)

var colorError *color.Color
var colorInfo *color.Color

var packetengineClient *packetengine.PacketEngineClient

var withoutTags string

func main() {
    colorError = color.New(color.FgWhite).Add(color.BgRed)
    colorInfo = color.New(color.FgGreen)

    app := &cli.App{
        Name:     "PacketEngine CLI",
        Usage:    "CLI application for the PacketEngine (packetengine.co.uk) API",
        Suggest:  true,
        Commands: []*cli.Command{
            {
                Name:        "init",
                HelpName:    "init",
                Action:      InitAction,
                ArgsUsage:   ` `,
                Usage:       `Sets the application API key.`,
                Description: `Set the API key.`,
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
        return errors.New("An API key is required.")
    }

    key := c.Args().First()

    err := SetAPIKey(key)
    if err != nil {
        return err
    }

    colorInfo.Println("API key set!")
    return nil
}

func SubdomainsAction(c *cli.Context) error {
    if c.Args().Len() < 1 {
        return errors.New("A domain is required.")
    }

    // Boot
    apiKey, err := GetAPIKey()

    if err != nil {
        return err
    }

    packetengineClient, err = packetengine.NewPacketEngineClient(apiKey)

    if err != nil {
        return err
    }

    domain := c.Args().First()

    subdomains, err := packetengineClient.GetSubdomains(domain, withoutTags)
    if err != nil {
        return err
    }

    for _, subdomain := range subdomains {
        colorInfo.Println(subdomain)
    }
    
    return nil
}

func SetAPIKey(apiKey string) error {
    if len(apiKey) == 0 {
        return errors.New("An API key is required.")
    }

    home, err := os.UserHomeDir()

    if err != nil {
        return err
    }

    keyBytes := []byte(apiKey)
    err = os.MkdirAll(home+"/.config/packetengine/", os.ModePerm)

    if err != nil {
        return err
    }

    err = os.WriteFile(home+"/.config/packetengine/api_key", keyBytes, 0644)

    if err != nil {
        return err
    }

    return nil
}

func GetAPIKey() (string, error) {
    home, err := os.UserHomeDir()

    if err != nil {
        return "", err
    }

    apiKey, err := os.ReadFile(home + "/.config/packetengine/api_key")

    if err != nil {
        return "", err
    }

    return string(apiKey), nil
}
