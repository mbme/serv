package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"io/ioutil"

	"path/filepath"

	"sync"

	"github.com/codegangsta/cli"
)

func readAsset(name string) []byte {
	data, err := Asset(name)

	if err != nil {
		panic(err)
	}

	return data
}

func readFile(name string) []byte {
	data, err := ioutil.ReadFile(name)

	if err != nil {
		panic(err)
	}

	return data
}

func servHTTP(wg *sync.WaitGroup, port string) {
	defer wg.Done()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func servHTTPS(wg *sync.WaitGroup, sslPort string, cert, key []byte) {
	defer wg.Done()

	srv := &http.Server{Addr: ":" + sslPort}

	config := &tls.Config{}
	config.NextProtos = []string{"http/1.1"}

	config.Certificates = make([]tls.Certificate, 1)

	var err error
	config.Certificates[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsListener, err := tls.Listen("tcp", srv.Addr, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := srv.Serve(tlsListener); err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "serv"
	app.Usage = "Simple static web server with SSL support"
	app.Version = "0.0.2"
	app.Author = "github.com/mbme"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 80,
			Usage: "http port",
		},
		cli.StringFlag{
			Name:  "dir,d",
			Value: ".",
			Usage: "directory to serv",
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "enable https",
		},
		cli.IntFlag{
			Name:  "ssl-port",
			Value: 443,
			Usage: "https port",
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
		dir, err := filepath.Abs(c.String("dir"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("serving directory", dir)
		http.Handle("/", http.FileServer(http.Dir(dir)))

		var wg sync.WaitGroup

		port := c.String("port")
		log.Println(" http: listening on port", port)
		wg.Add(1)
		go servHTTP(&wg, port)

		if c.Bool("ssl") {
			sslPort := c.String("ssl-port")
			log.Printf("https: listening on port %v", sslPort)

			certPath := c.String("cert")
			keyPath := c.String("key")

			var cert []byte
			var key []byte
			if len(certPath) == 0 || len(keyPath) == 0 {
				log.Println("https: using embedded certificate and key")

				cert = readAsset("server.crt")
				key = readAsset("server.key")
			} else {
				log.Println("https: using provided certificate and key")

				cert = readFile(certPath)
				key = readFile(keyPath)
			}

			wg.Add(1)
			go servHTTPS(&wg, sslPort, cert, key)
		}

		wg.Wait()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
