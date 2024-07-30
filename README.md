# Ecowitt Proxy

This is a small command-line utility written in Go to proxy Ecowitt weatherstation updates to Home Assistant.

## Usage

```
export PROXY_PORT=4199
export DEST_URL=https://my.homeassistant.com/api/webhook/123312312897345987239847
./ecoproxy
```

## Building

To build the CLI tool, you can either use your local golang installation or have Docker build it.

### Docker

Run the following command to build the Docker image:

```
$ docker build -t ecoproxy:latest .
```

After the build is complete, you can run the CLI tool using the Docker image:

```
$ docker run --rm -e BASE_URL=<endpoint> ecoproxy:latest
```

This repo also has an example `docker-compose.yaml` that can be used.
Works for me (tm).

### Local build

Run the following command to build the chat binary:
```
$ make build-local
```

## TODO/Notes

Nothing critical right now.
Options when bored:
 - support for multiple weatherstation/endpoint combos (just run it multiple times for now)
