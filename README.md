# `qc`

[![License](https://img.shields.io/github/license/troykinsella/qc.svg)](https://github.com/troykinsella/qc/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/qc.svg?branch=master)](https://travis-ci.org/troykinsella/qc)
[![Go Report](https://goreportcard.com/badge/github.com/troykinsella/qc)](https://goreportcard.com/report/github.com/troykinsella/qc)

> "queue cat"

A command line utility for publishing and subscribing to 
AMQP 0.9.1 message queues, like RabbitMQ.

Note: This application is in alpha status until a `1.0.0` release.

## Installation

Head over to [releases](https://github.com/troykinsella/qc/releases) and download the appropriate binary for your system.
Put the binary in a convenient place, such as `/usr/local/bin/qc`.

## Usage

By default, `qc` connects to a message queue broker at
`amqp://guest:guest@localhost:5672/`. Pass `-u` with 
an alternate URL, or set the `QC_BROKER_URL` environment variable
to change the broker to which `qc` connects.

Supply an exchange name with the `-e` option.
Supply one or more routing keys with the `-r` option.
Pass `-t` with value `direct` (default), `fanout`, `headers`, or `topic` 
to control the exchange type.

`qc` can be executed in producer or consumer mode.

### Consumer Mode

Consume messages by passing the `-c` option, 
an exchange name (`-e`), and some routing keys (`-r`) 
on which to bind.

```bash
qc -c -e logs -r info -r error
```

This declares an exchange named `logs` of the default
type `direct`. It creates an exclusive queue having a 
generated name, and binds that queue to the exchange for 
each routing key supplied (`info`, and `error`), or once
with an empty routing key if none is supplied. Lastly, 
it consumes messages from the queue, writing them to `stdout`, 
until the process is killed.

### Producer Mode

Producer mode is the default, active when `-c` is omitted.

Publish messages by passing input to `qc` through `stdin`,
supplying the exchange name (`-e`) and a routing key (`-r`).

```bash
echo "oh no!" | qc -e logs -r error
```

This declares an exchange named `logs` of the default
type `direct`, and publishes the given message to it.

## Examples

### Direct Exchange

Create separate consumers that are interested only in specific
log levels on exchange `logs`. As `direct` is the default 
exchange type, we can omit the `-t` option.

```bash
# error handler
$ qc -c -e logs -r error
```

```bash
# info handler
$ qc -c -e logs -r info
```

Now log some messages.

```bash
# logger A
$ echo "your bike was stolen!" | qc -e logs -r error
$ echo "your mother-in-law was kidnapped" | qc -e logs -r info
```

Consumers receive only the relevant messages.

```bash
# error handler
$ qc -c -e logs -r error
your bike was stolen!
```

```bash
# info handler
$ qc -c -e logs -r info
your mother-in-law was kidnapped
```

### Pub-Sub / Fan-Out

Create some consumers to receive log messages on exchange `logs`.
Routing keys are not needed for `fanout`.

```bash
# consumer A
$ qc -c -t fanout -e logs
```

```bash
# consumer B
$ qc -c -t fanout -e logs
```

Publish log messages.

```bash
# producer A
$ echo "Houston, there's cake" | qc -t fanout -e logs
```

Both consumers receive the message.

```bash
# consumer A
$ qc -c -t fanout -e logs
Houston, there's cake
```

```bash
# consumer B
$ qc -c -t fanout -e logs
Houston, there's cake
```

### Topics

Create consumers that are interested in particular topics.

```bash
# any-colour rabbit consumer
qc -c -t topic -e animals -r '*.rabbit'
```

```bash
# blue animal consumer
qc -c -t topic -e animals -r 'blue.*'
```

Publish some animal names.

```bash
# producer A
echo "roger" | qc -t topic -e animals -r white.rabbit
echo "dick"  | qc -t topic -e animals -r blue.whale
echo "ted"   | qc -t topic -e animals -r blue.rabbit
echo "jerry" | qc -t topic -e animals -r red.eagle
```

Consumers receive only the relevant messages.

```bash
# any-colour rabbit consumer
qc -c -t topic -e animals -r '*.rabbit'
roger
ted
```

```bash
# blue animal consumer
qc -c -t topic -e animals -r 'blue.*'
dick
ted
```

## License

MIT © Troy Kinsella
