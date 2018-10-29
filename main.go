package main

import (
	"fmt"
	"log"
	"os"

	"github.com/samze/broker-watch/broker"
	"github.com/urfave/cli"
)

func main() {
	var (
		url      string
		username string
		password string
	)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "url",
			Usage:       "broker url",
			Destination: &url,
		},
		cli.StringFlag{
			Name:        "username",
			Usage:       "broker username",
			Destination: &username,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "broker password",
			Destination: &password,
		},
	}

	app.Name = "broker watcher"
	app.Usage = "check broker availability"
	app.Action = func(c *cli.Context) error {

		if url == "" || username == "" || password == "" {
			log.Fatal("Missing url, username or password")
		}

		broker, err := broker.NewBroker(url, username, password)
		if err != nil {
			log.Fatalf("error creating broker client: +%v", err)
		}

		err = broker.Check()
		if err != nil {
			log.Fatalf("broker failed check +%v", err)
		}

		fmt.Println("broker passed check")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
