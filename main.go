package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/aurelius15/product-reviews/cmd/api"
	"github.com/aurelius15/product-reviews/cmd/migration"
	"github.com/aurelius15/product-reviews/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cmd := &cli.Command{
		Name:    "product-reviews",
		Usage:   "This service is responsible for providing product ratings",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db-host",
				Usage:    "DB host to connect (example: 'localhost:5432')",
				Sources:  cli.EnvVars("DB_HOST"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-name",
				Usage:    "DB name to connect",
				Sources:  cli.EnvVars("DB_NAME"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-user",
				Usage:    "DB user to connect",
				Sources:  cli.EnvVars("DB_USER"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "db-password",
				Usage:    "DB password to connect",
				Sources:  cli.EnvVars("DB_PASSWORD"),
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Launch API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "nats-host",
						Usage:    "NATS host to connect (example: 'localhost:4222')",
						Sources:  cli.EnvVars("NATS_HOST"),
						Required: true,
					},
					&cli.StringFlag{
						Name:     "nats-subject",
						Usage:    "NATS subject to listen",
						Sources:  cli.EnvVars("NATS_SUBJECT"),
						Required: true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return api.RestAPICmd(ctx, config.NewPostgresCnf(c), config.NewNATSCnf(c))
				},
			},
			{
				Name:  "migration",
				Usage: "Run DB migrations",
				Action: func(_ context.Context, c *cli.Command) error {
					return migration.RunMigrationCmd(config.NewPostgresCnf(c))
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("service exit with error", slog.Any("err", err))
		os.Exit(1)
	}
}
