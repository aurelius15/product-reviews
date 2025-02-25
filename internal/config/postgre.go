package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

type PostgresCnf struct {
	host     string
	port     int
	db       string
	user     string
	password string
}

func NewPostgresCnf(c *cli.Command) *PostgresCnf {
	hostAndPort := strings.Split(c.String("db-host"), ":")
	port, _ := strconv.Atoi(hostAndPort[1])

	return &PostgresCnf{
		host:     hostAndPort[0],
		port:     port,
		db:       c.String("db-name"),
		user:     c.String("db-user"),
		password: c.String("db-password"),
	}
}

func (p *PostgresCnf) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=public",
		p.user, p.password, p.host, p.port, p.db)
}

func (_ *PostgresCnf) PathToMigrations() string {
	return "./migrations"
}
