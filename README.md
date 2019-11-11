# circle-dd-bench

Wraps an arbitrary command run in CircleCI and send the running time to Datadog.

## Usage

```
$ circle-dd-bench [OPTIONS] -- COMMAND
```

### Options

* `-t`, `--tag=`: Tag to send to Datadog with TAG:VALUE format

### Environment Variables

* `DATADOG_API_KEY`: Datadog API key

## Installation

### Docker

```
# In Dockerfile
COPY --from=yuyat/circle-dd-bench /usr/local/bin/circle-dd-bench /usr/local/bin
```

### Build yorself

```
$ go get github.com/yuya-takeyama/circle-dd-bench
```
