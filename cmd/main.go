package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"alertbot/config"
	"alertbot/internal/consumer"
	"alertbot/internal/events/telegram"
	"alertbot/internal/usecase"
)

const (
	batchSize = 100
)

func main() {
	application := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: true,
				Usage:    "YAML config filepath",
				EnvVars:  []string{"CONFIG_FILE"},
				FilePath: "/srv/secret/config_file",
			},
		},
		Action: Main,
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Main(ctx *cli.Context) error {
	cfg, err := config.New(ctx.String("config-file"))
	if err != nil {
		return err
	}

	eventProcessor := telegram.NewTg(usecase.New(cfg.AlertBot.Host, cfg.AlertBot.Token))

	cons := consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := cons.Start(); err != nil {
		return err
	}

	return nil
}
