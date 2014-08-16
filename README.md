# Proxyd

A distributed fault tolerant HTTP proxy.

## Getting started

### Building

You can build Proxyd from source

```sh
git clone git@github.com:AlexisMontagne/proxyd.git
cd proxyd
make
```

This will generate to binaries `./build/proxyd_endpoint` and
`./build/proxyd_balancer`

_NOTE_: you need a Go version to compile it.

### Running

You can launch multiple endpoints easily with

```sh
./build/proxyd_endpoint -port 1080
./build/proxyd_endpoint -port 1081
...
```


Then launch a balancer to distribute the HTTP calls to the endpoints

```sh
./build/proxyd_balcancer -port 1079
```

Then you can use proxyd to proxy your request

```sh
curl --proxy localhost:1079 http://www.google.com
```
