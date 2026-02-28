package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ngyewch/go-opensky"
	"github.com/urfave/cli/v3"
)

const (
	sgMinLat = 1.129
	sgMinLon = 103.557
	sgMaxLat = 1.493
	sgMaxLon = 104.131
)

func doRetrieve(ctx context.Context, cmd *cli.Command) error {
	extended := cmd.Bool(extendedFlag.Name)
	minLat := cmd.Float64(minLatFlag.Name)
	minLon := cmd.Float64(minLonFlag.Name)
	maxLat := cmd.Float64(maxLatFlag.Name)
	maxLon := cmd.Float64(maxLonFlag.Name)

	client := newClient(ctx, cmd)

	allStatesResponse, err := client.GetAllStates(opensky.AllStatesRequest{
		Extended: extended,
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
