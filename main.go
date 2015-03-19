package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "servs"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "websockets port",
		},
		cli.StringFlag{
			Name:  "cert",
			Value: "",
			Usage: "ssl certificate path",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "",
			Usage: "ssl certificate key path",
		},
	}

	app.Action = func(c *cli.Context) {
		port := c.String("port")
		log.Printf("listening on port %v", port)

		http.Handle("/", http.FileServer(http.Dir(".")))

		certPath := c.String("cert")
		keyPath := c.String("key")

		if len(certPath) == 0 || len(keyPath) == 0 { //  HTTP
			log.Println("http-only mode")
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				log.Fatal(err)
			}
		} else { // HTTPS
			log.Println("enabled SSL")
			if err := http.ListenAndServeTLS(":"+port, certPath, keyPath, nil); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
