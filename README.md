# serv
static file server with SSL support

## DEPENDENCIES

* go 1.5+ to compile app

* GNU make to generate new certificates
* openssl to generate new certificates
* [go-bindata](https://github.com/jteeuwen/go-bindata) to embed new certificates into app

## TODO

* serv both http and https
* avoid using openssl, use internal crypto/tls package to generate keys
