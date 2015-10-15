# serv
static file server with SSL support

# DEPENDENCIES

* go 1.5+ to compile app
* [gb](https://github.com/constabulary/gb) as go build tool

* GNU make to generate new certificates
* openssl to generate new certificates
* [go-bindata](https://github.com/jteeuwen/go-bindata) to embed new certificates into app

# BUILD

Run `make build` to compile app.

To generate new certificates run `make gen-certs build`.

# USE

run `serv -h` to see help.

# TODO

* avoid using openssl, use internal crypto/tls package to generate keys
