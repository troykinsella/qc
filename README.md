# `qc`

[![License](https://img.shields.io/github/license/troykinsella/qc.svg)](https://github.com/troykinsella/qc/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/qc.svg?branch=master)](https://travis-ci.org/troykinsella/qc)

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
an alternate URL to change the broker to which `qc` connects.

`qc` can be executed in producer or consumer mode.

### Consumers

Consume messages by passing the `-c` option, 
an exchange name, and some routing keys on which to bind:

```bash
qc -c logs info warning error
```

This declares an exchange named `logs` of the default
type `direct`. You can pass `-t` with value `direct`, 
`fanout`, or `topic` to control the exchange type. It creates 
an exclusive queue having a generated name, and binds it to 
the exchange for each routing key supplied 
(`info`, `warning`, and `error`). Lastly, it consumes
messages from the queue, writing them to `stdout`, until the
process is killed.

### Producers

Produce messages by passing input to `qc` through `stdin`,
supplying the exchange name and a routing key:

```bash
echo "oh no!" | qc logs error
```

This declares an exchange named `logs`, and publishes the
given message to it.

## Roadmap

* More control over exchange and queue options

## License

MIT © Troy Kinsella
