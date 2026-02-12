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

	clientIdFlag = cli.StringFlag{
		Name:    "client-id",
		Usage:   "client ID",
		Sources: cli.EnvVars("CLIENT_ID"),
	}
	clientSecretFlag = cli.StringFlag{
		Name:    "client-secret",
		Usage:   "client secret",
		Sources: cli.EnvVars("CLIENT_SECRET"),
	}

	app = &cli.Command{
		Name:    "opensky",
		Usage:   "OpenSky CLI",
		Version: version,
		Action:  nil,
		Commands: []*cli.Command{
			{
				Name:   "test",
				Usage:  "test",
				Action: doTest,
				Flags: []cli.Flag{
					&clientIdFlag,
					&clientSecretFlag,
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
