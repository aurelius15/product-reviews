package config

import (
	"github.com/urfave/cli/v3"
)

type RedisCnf struct {
	host string
}

func NewRedisCnf(c *cli.Command) *RedisCnf {
	return &RedisCnf{
		host: c.String("redis-host"),
	}
}

func (r *RedisCnf) Host() string {
	return r.host
}
