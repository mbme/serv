package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"io/ioutil"

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

func listenAndServ(srv *http.Server, cert, key []byte) error {
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}

	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)

	var err error
	config.Certificates[0], err = tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}

	tlsListener, err := tls.Listen("tcp", srv.Addr, config)
	if err != nil {
		return err
	}

	return srv.Serve(tlsListener)
}

func main() {
	app := cli.NewApp()
	app.Name = "servs"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 8080,
			Usage: "websockets port",
		},
		cli.BoolFlag{
			Name:  "ssl",
			Usage: "enable https",
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

		if !c.Bool("ssl") {
			log.Println("SSL disabled")
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				log.Fatal(err)
			}
			return
		}

		log.Println("SSL enabled")

		certPath := c.String("cert")
		keyPath := c.String("key")

		var cert []byte
		var key []byte
		if len(certPath) == 0 || len(keyPath) == 0 {
			log.Println("using embedded certificate and key")

			cert = readAsset("server.crt")
			key = readAsset("server.key")
		} else {
			log.Println("using provided certificate and key")

			cert = readFile(certPath)
			key = readFile(keyPath)
		}

		server := &http.Server{Addr: ":" + port}
		if err := listenAndServ(server, cert, key); err != nil {
			log.Fatal(err)
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
