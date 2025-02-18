package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	DSN string `mapstructure:"NATS_DSN"`
}

func New(cfg *Config) (jetstream.JetStream, func(), error) {
	nc, err := nats.Connect(cfg.DSN)
	if err != nil {
		return nil, nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, err
	}

	cleaner := func() {
		nc.Close()
	}

	return js, cleaner, nil
}
