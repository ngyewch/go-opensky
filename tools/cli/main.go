package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
)

var (
	version string

	clientIdFlag = &cli.StringFlag{
		Name:    "client-id",
		Usage:   "client ID",
		Sources: cli.EnvVars("CLIENT_ID"),
	}
	clientSecretFlag = &cli.StringFlag{
		Name:    "client-secret",
		Usage:   "client secret",
		Sources: cli.EnvVars("CLIENT_SECRET"),
	}
	extendedFlag = &cli.BoolFlag{
		Name:    "extended",
		Usage:   "extended",
		Sources: cli.EnvVars("EXTENDED"),
	}
	minLatFlag = &cli.Float64Flag{
		Name:     "min-lat",
		Usage:    "minimum latitude",
		Required: true,
		Sources:  cli.EnvVars("MIN_LAT"),
	}
	minLonFlag = &cli.Float64Flag{
		Name:     "min-lon",
		Usage:    "minimum longitude",
		Required: true,
		Sources:  cli.EnvVars("MIN_LON"),
	}
	maxLatFlag = &cli.Float64Flag{
		Name:     "max-lat",
		Usage:    "maximum latitude",
		Required: true,
		Sources:  cli.EnvVars("MAX_LAT"),
	}
	maxLonFlag = &cli.Float64Flag{
		Name:     "max-lon",
		Usage:    "maximum longitude",
		Required: true,
		Sources:  cli.EnvVars("MAX_LON"),
	}

	app = &cli.Command{
		Name:    "opensky",
		Usage:   "OpenSky CLI",
		Version: version,
		Action:  nil,
		Commands: []*cli.Command{
			{
				Name:   "retrieve",
				Usage:  "retrieve",
				Action: doRetrieve,
				Flags: []cli.Flag{
					clientIdFlag,
					clientSecretFlag,
					extendedFlag,
					minLatFlag,
					minLonFlag,
					maxLatFlag,
					maxLonFlag,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		oopsError, ok := oops.AsOops(err)
		if ok {
			fmt.Printf("%+v", oopsError)
		} else {
			log.Fatal(err)
		}
	}
}
