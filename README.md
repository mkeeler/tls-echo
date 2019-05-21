# tls-echo

`tls-echo` is a simple TLS enabled echo server.

## Installation

`go get github.com/mkeeler/tls-echo`

## Usage

```
> tls-echo -cert <path to server certificate> -key <path to private key> -listen ":443" -prefix "echo prefix: "
2019/05/21 13:37:39 Now accepting TLS connections for :443
```
