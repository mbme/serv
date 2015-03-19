CERTS := ./certs

clean-certs:
	rm -rf $(CERTS)
	mkdir $(CERTS)

gen-certs: clean-certs
	openssl genrsa -des3 -passout pass:x -out $(CERTS)/server.pass.key 2048
	openssl rsa -passin pass:x -in $(CERTS)/server.pass.key -out $(CERTS)/server.key
	rm $(CERTS)/server.pass.key

	openssl req -new -key $(CERTS)/server.key -out $(CERTS)/server.csr
	openssl x509 -req -days 365 -in $(CERTS)/server.csr -signkey $(CERTS)/server.key -out $(CERTS)/server.crt
	go-bindata -o certs.go -nomemcopy -prefix "certs" certs/...

deps:
	go get -u github.com/codegangsta/cli
