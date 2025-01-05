package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	Version = "dev"
)

func main() {
	app := &cli.App{
		Name:  "spirrel",
		Usage: "CLI for generating and storing articles in Elasticsearch",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "esUri",
						Aliases: []string{"s"},
						Usage:   "Elasticsearch host",
						EnvVars: []string{"ELASTICSEARCH_HOST"},
						Value:   "http://localhost:9200",
					},
					&cli.StringFlag{
						Name:     "esApiKey",
						Aliases:  []string{"k"},
						Usage:    "Elasticsearch API Key",
						EnvVars:  []string{"ELASTICSEARCH_API_KEY"},
						Required: true,
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "Port to run the API server on",
						Value:   8080,
					},
				},
				Action: func(c *cli.Context) error {
					esURL := c.String("esUri")
					port := c.Int("port")
					esApiKey := c.String("esApiKey")
					return runServer(esURL, esApiKey, port)
				},
			},
		},
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("")
	}
}
