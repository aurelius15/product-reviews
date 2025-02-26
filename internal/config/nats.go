package config

import (
	"github.com/urfave/cli/v3"
)

type NatsCnf struct {
	host    string
	subject string
}

func NewNATSCnf(c *cli.Command) *NatsCnf {
	return &NatsCnf{
		host:    c.String("nats-host"),
		subject: c.String("nats-subject"),
	}
}

func (c NatsCnf) Host() string {
	return c.host
}

func (c NatsCnf) Subject() string {
	return c.subject
}
