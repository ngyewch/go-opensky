package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func newClient(ctx context.Context, cmd *cli.Command) *Client {
	clientId := cmd.String(clientIdFlag.Name)
	clientSecret := cmd.String(clientSecretFlag.Name)

	return NewClient(clientId, clientSecret, nil)
}

func doTest(ctx context.Context, cmd *cli.Command) error {
	client := newClient(ctx, cmd)

	const (
		minLat = 1.129
		minLon = 103.557
		maxLat = 1.493
		maxLon = 104.131
	)

	allStatesResponse, err := client.GetAllStates(AllStatesRequest{
		Extended: true,
		MinLat:   minLat,
		MinLon:   minLon,
		MaxLat:   maxLat,
		MaxLon:   maxLon,
	})
	if err != nil {
		return err
	}

	err = os.MkdirAll("output", 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join("output", "output.csv"))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	err = allStatesResponse.WriteToCSV(f)
	if err != nil {
		return err
	}

	return nil
}
