# Mutual certificate authentication using HAProxy sidecars

This repo shows how to use HAProxy as a sidecar that runs alongside your main
application. The sidecar HAProxy implements mutual TLS authentication using
certificates that share a Certficate Authority. Both client and server sides
of the HAProxy connection present certificates for verification.

Broadly this shows how HAProxy may be used as a proxy to upgrade the connection
between two applications to a TLS connection with mutual authentication. In
particular this is useful to tunnel otherwise non-TLS capable clients over a
secured, authenticated channel.

The main files are:

- `client-proxy.cfg` shows a HAProxy configuration using a `server` line which
    creates a TLS connection to the server proxy, presenting a certificate
    for authentication.
- `server-proxy.cfg` shows a HAProxy configuration using a `bind` line which
    receives a TLS connection to the client proxy, presenting a certificate
    for authentication.
- `main.go` contains a simple client and server to demonstrate the proxy in use.
- `Makefile` and `openssl-with-ca.conf` show the relevant `openssl` commands
    to generate a Certificate Authority (CA) signing key and certificate, along
    with the signing keys and certificates for the client and server HAProxy
    instances to use.

The (imagined) setup is:

```
┌───────────────────────────────┐           ┌──────────────────────────────┐
│                               │           │                              │
│ ┌────────┐      ┌────────────┬┤           ├┬────────────┐      ┌───────┐ │
│ │ go     │:8080 │            ││  :8081    ││            │:3000 │go     │ │
│ │ client ├─────►│client-proxy│┼───────────┼│server-proxy├─────►│server │ │
│ │        │      │            ││  secure   ││            │      │       │ │
│ └────────┘      └────────────┴┤  TLS 1.3+ ├┴────────────┘      └───────┘ │
│                               │  channel  │                              │
│           node 1              │           │            node 2            │
└───────────────────────────────┘           └──────────────────────────────┘
```

## Getting started

First, we need to generate the certificates needed for the HAProxy channel:

- `ca.crt` is the Certificate Authority's certificate, used to verify the
    connection between the client and server HAProxy.
- `client.pem` has a certificate and signing key which has been signed by the
    Certificate Authority's signing key for use by the client HAProxy.
- `server.pem` has a certificate and signing key which has been signed by the
    Certificate Authority's signing key for use by the server HAProxy.

Generate these using:

```
make all-certs
```

All of these files will be written to `./certs/`, along with the Certificate
Authority's signing key, `ca.key`.

## Running the demo

First, start the Go server:

```
go run . --mode server --port 3000
```

Next, start the server HAProxy:

```
haproxy -f server-proxy.cfg
```

Next, the client HAProxy:

```
haproxy -f client-proxy.cfg
```

Neither of the `haproxy` invocations should print output.

Finally, run the Go client application:

```
go run . --mode client --port 8080
