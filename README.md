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

#### Barebone check

If you want to write a simple check you just have to write a function
with zero argument and which returns a check.ExtensionCheckResult. Wrap
this function pointer into an `check.ExtensionCheck` and add it into the
`check.Store` map with check name as key,  such as:

package main

```golang
package main

import (
	"net/http"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/handler"
)

func HTTPCheck() check.ExtensionCheckResult {
	resp, err := http.Get("http://example.com/")

	if err != nil {
		return handler.Error(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return handler.Ok("Example.com is alive!")
	} else if resp.StatusCode >= 500 {
		return handler.Error("Example.com return an 5XX status code")
	} else {
		return handler.Warning("Example.com is not responding as expected")
	}
}

func main() {
	cfg := sensu.NewConfigFromFlagSet(sensu.ExtractFlags())

	t := rabbitmq.NewRabbitMQTransport(cfg.RabbitMQURI())
	client := sensu.NewClient(t, cfg)

	check.Store["http_check"] = &check.ExtensionCheck{HTTPCheck}

	client.Start()
}
```

#### Standard check

If you want to write a check and a metric which inspect the same value
you can use the `StandardCheck` struct from the `github.com/upfluence/sensu-client-go/sensu/handler`
package. Such as:

```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/upfluence/sensu-client-go/Godeps/_workspace/src/github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/utils"
)

func HTTPCallDuration() (float64, error) {
	t0 := time.Now().Unix()
	_, err := http.Get("http://example.com/")

	return float64(time.Now().Unix() - t0), err
}

func main() {
	c := &utils.StandardCheck{
		ErrorThreshold:   20.0,
		WarningThreshold: 10.0,
		MetricName:       "http_call.duration",
		Value:            HTTPCallDuration,
		CheckMessage: func(v float64) string {
			return fmt.Sprintf("Duration: %.2fs", v)
		},
		Comp: func(x, y float64) bool { return x > y },
	}

	cfg := sensu.NewConfigFromFlagSet(sensu.ExtractFlags())

	t := rabbitmq.NewRabbitMQTransport(cfg.RabbitMQURI())
	client := sensu.NewClient(t, cfg)

	check.Store["http_duration_check"] = &check.ExtensionCheck{c.Check}
	check.Store["http_duration_metric"] = &check.ExtensionCheck{c.Metric}

	client.Start()
}
```

If the duration exceed 20s the `http_duration_check` will return an
Error, if the duration exceed 10s the check will return a warning
otherwise it will returns an "OK" code.

The `http_duration_metric` will return `http_call_duration 5.000000 1438125085`

with 5.0 as the duration of the HTTP call and 1438125085 the timestamp

### Running

You just have to compile it and execute it, such as:

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
| SENSU_CLIENT_SUBSCRIPTIONS | Comma separated subscriptions | email,slack |
| SENSU_CLIENT_NAME | The name of the client | node-01 |
| SENSU_CLIENT_ADDRESS | The ip addres of the client | 127.0.0.1 |
| RABBITMQ_URL | RabbitMQ url | amqp://guest:guest@localhost:5672/%2f |

## Roadmap

* [ ] Implement the keep-alives specific configurations (thresholds and
  hanlder)
* [ ] Implement the standalone check mechanism
