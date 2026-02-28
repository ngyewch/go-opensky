package main

import (
	"context"

	"github.com/ngyewch/go-opensky"
	"github.com/urfave/cli/v3"
)

func newClient(ctx context.Context, cmd *cli.Command) *opensky.Client {
	clientId := cmd.String(clientIdFlag.Name)
	clientSecret := cmd.String(clientSecretFlag.Name)

	return opensky.NewClient(clientId, clientSecret, nil)
}
