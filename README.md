# Mutual certificate authentication using HAProxy sidecars

This repository is meant to demonstrate the concept of using a proxy to create a
secure channel. It shows how to use HAProxy to create a secure (encrypted and
authenticated) channel between two applications that don't support secure
communication themselves. It uses HAProxy as a "sidecar" at each end of a
client/server connection. The client connects to its local HAProxy, which
connects to the server's HAProxy sidecar, which connects to the server.

The client -> HAProxy and HAProxy -> server connections are unencrypted as they
are assumed to be machine-local. The client and server HAProxy instances use
a shared Certificate Authority (CA) to trust certificates that they present
to each other when making the connection.

As a reference point, this is essentially what tools like Istio and Linkerd
do to create secure channels between pods in Kubernetes -- create a proxy which
has its own certificate signed by a shared CA. In addition, other proxies like
Envoy and Nginx can easily be used in place of HAProxy; I just happen to know
HAProxy best.

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
```

## Troubleshooting

Errors when starting HAProxy that contain `DOWN`:

```
[WARNING]  (73233) : Server client/s1 is DOWN, reason: Layer4 connection problem, info: "Connection refused", check duration: 0ms. 0 active and 0 backup servers left. 0 sessions active, 0 requeued, 0 remaining in queue.
```

What this means:

- From the client HAProxy, this means you started the client HAProxy before
    starting the server HAProxy.
- From the server HAProxy, this means you started the server HAProxy before
    starting the Go server application.
