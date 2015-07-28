# Sensu Client
  Portable / Embeddable lighweight version of the sensu client

[![Circle CI](https://circleci.com/gh/upfluence/sensu-client-go.svg?style=svg)](https://circleci.com/gh/upfluence/sensu-client-go)

## Usage

You can use this project by two different ways.

  * The first is use the project as a standalone program and use
    external checks/metrics.

  * The second is to use the project as an external golang package and
    write a project with another main method and add some checks /
    metrics inside.


### Standalone project

The easiest way to get started is to download the binary from your
command line:

* Linux

```shell
$ curl -sL https://github.com/upfluence/sensu-client-go/releases/download/v0.0.1/sensu-client-go-linux-amd64-0.0.1 \
  > sensu-client

$ chmod +x sensu-client
```

* OSX

```shell
$ curl -sL
https://github.com/upfluence/sensu-client-go/releases/download/v0.0.1/sensu-client-go-darwin-amd64-0.0.1 \
  > sensu-client

$ chmod +x sensu-client
```

Then you just have to run the program as `./sensu-client -c /etc/sensu/config.json`. EASY!

### Subpackage

You just have to use the tools packaged into this repository. Hereby an
example:

```golang
package main

import (
  "net/http"
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/transport"
	"github.com/upfluence/sensu-client-go/sensu/handler"
)

func HTTPCheck() check.ExtensionCheckResult {
  resp, err := http.Get("http://example.com/")

  if err != nil {
    return handler.Error(err.Error())
  }

  if resp.Status >= 200 && resp.Status < 300 {
    return handler.Ok("Example.com is alive!")
  } else if resp.Status >= 500 {
    return handler.Error("Example.com return an 5XX status code")
  } else {
    return handler.Warning("Example.com is not responding as expected")
  }
}

func main() {
	cfg := sensu.NewConfigFromFlagSet(sensu.ExtractFlags())

	t := transport.NewRabbitMQTransport(cfg)
	client := sensu.NewClient(t, cfg)

	check.Store["http_check"] = &check.ExtensionCheck{HTTPCheck}

	client.Start()
}
```

Then, just compile it and run it, such as:

```shell
$ ls

my_awesome_main.go sensu-config.json

$ go build -o sensu-client .

$ ./sensu-client -c sensu-config.json

```

### Options

In the both cases, you can use the  `-c` flag to use a specific
configuration file, the configuration is pretty similar to the ruby
client [check out the doc](http://sensuapp.org/docs/0.16/clients). The
difference is about the configuration of the RabbitMQ client. You have
to provide an RabbitMQ URI through the `RABBITMQ_URL` environment
variable or by adding the `rabbit_uri` key into the root of the JSON
configuration file.

By the way you can also specify some options through environment
variables:

| variable | explanation | example |
| -------  | ----------- | ------  |
| SENSU_CLIENT_SUVSCRIPTIONS | Comma separated subscriptions | email,slack |
| SENSU_CLIENT_NAME | The name of the client | node-01 |
| SENSU_CLIENT_ADDRESS | The ip addres of the client | 127.0.0.1 |
| RABBITMQ_URL | RabbitMQ url | amqp://guest:guest@localhost:5672/%2f |

## Roadmap

* [ ] Implement some metrics check helper
* [ ] Implement the keep-alives specific configurations (thresholds and
  hanlder)
* [ ] Implement the standalone check mechanism
