.PHONY: cert-dir ca-cert server-cert client-cert

cert-dir:
	mkdir -p certs

all-certs: server-cert client-cert

# Create our own trusted CA. Well, HAProxy will trust it in the configuration
# files in this repository anyway.
ca-cert: cert-dir
	# openssl genrsa -out ./certs/ca.key 4096
	openssl req \
		-config ./openssl-with-ca.cnf \
		-new -nodes -x509 -days 36135 -extensions v3_ca -newkey rsa:4096 \
		-subj "/C=/ST=/L=/O=/CN=foo-ca.dx13.co.uk" \
		-keyout ./certs/ca.key -out ./certs/ca.crt

# Create certificate for proxy-server.cfg which is signed by the CA
server-cert: cert-dir ca-cert
	openssl genrsa -out ./certs/server.key 4096
	openssl req -new \
		-key ./certs/server.key -out ./certs/server.csr \
		-subj "/C=/ST=/L=/O=/CN=foo-server.dx13.co.uk"
	openssl x509 -req -days 36135 \
		-in ./certs/server.csr \
		-CA ./certs/ca.crt -CAkey ./certs/ca.key \
		-set_serial 01 \
		-out ./certs/server.crt
	cat ./certs/server.key > ./certs/server.pem
	cat ./certs/server.crt >> ./certs/server.pem
	rm ./certs/server.key ./certs/server.crt ./certs/server.csr

# Create certificate for proxy-client.cfg which is signed by the CA
client-cert: cert-dir ca-cert
	openssl genrsa -out ./certs/client.key 4096
	openssl req -new -key ./certs/client.key -out ./certs/client.csr \
		-subj "/C=/ST=/L=/O=/CN=foo-client.dx13.co.uk"
	openssl x509 -req -days 36135 \
		-in ./certs/client.csr \
		-CA ./certs/ca.crt -CAkey ./certs/ca.key \
		-set_serial 01 \
		-out ./certs/client.crt
	cat ./certs/client.key > ./certs/client.pem
	cat ./certs/client.crt >> ./certs/client.pem
	rm ./certs/client.key ./certs/client.crt ./certs/client.csr

clean:
	- rm -r certs
