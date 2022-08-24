package config

import (
	"fmt"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/tg/pgpass"
	nurl "net/url"
	"sort"
	"strings"
)

type Bastion struct {
	User string `yaml:"user,omitempty"`
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

type Ssh struct {
	User    string  `yaml:"user,omitempty"`
	Host    string  `yaml:"host,omitempty"`
	Port    int     `yaml:"port,omitempty"`
	Key     string  `yaml:"key,omitempty"`
	Bastion Bastion `yaml:"bastion,omitempty"`
}

type Connection struct {
	Name string `yaml:"name,omitempty"`
	URL  string `yaml:"url,omitempty"`
	Ssh  Ssh    `yaml:"ssh,omitempty"`

	connStr string
}

type Toolbox struct {
	Connections []Connection `yaml:"connections,omitempty"`
	PGDATA      string       `yaml:"pgdata,omitempty"`
	PGBIN       string       `yaml:"pgbin,omitempty"`
}

type Config struct {
	Toolbox Toolbox `yaml:"toolbox,omitempty"`
}

func (c *Config) GetPGDATA() string {
	return c.Toolbox.PGDATA
}

func (c *Config) GetPGBIN() string {
	return c.Toolbox.PGBIN
}

func (c *Config) GetConn(name string) *Connection {
	for _, conn := range c.Toolbox.Connections {
		if conn.Name == name {
			return &conn
		}
	}
	return nil
}

func (c *Config) GetAllConnNames() []string {
	names := make([]string, 0, len(c.Toolbox.Connections))
	for _, connection := range c.Toolbox.Connections {
		names = append(names, connection.Name)
	}
	sort.Strings(names)
	return names
}

func (c *Connection) ParseURL() string {
	url, err := pgpass.UpdateURL(c.URL)
	if err != nil {
		log.WithField("conn", c.Name).WithError(err).Fatal()
		return ""
	}
	if !checkPass(url) {
		log.WithField("conn", c.Name).Fatal("no password provided for user")
		return ""
	}

	connStr, err := pq.ParseURL(url)
	if err != nil {
		log.WithField("conn", c.Name).
			WithError(err).
			Fatalf("cannot parse %s url", c.Name)
		return ""
	}

	c.connStr = connStr

	return connStr
}
func checkPass(url string) bool {
	u, err := nurl.Parse(url)
	if err != nil {
		return false
	}
	if user := u.User; user != nil {
		if _, ok := user.Password(); !ok {
			return false
		}
	}
	return true
}

func (c *Connection) GetFieldFromConnStr(field string) (string, error) {
	if c.connStr == "" {
		c.ParseURL()
	}

	params := strings.Split(c.connStr, " ")
	for _, param := range params {
		if strings.Contains(param, field) {
			value := strings.Split(param, "=")[1]
			value = strings.Trim(value, "'")
			return value, nil
		}
	}
	return "", fmt.Errorf("cannot find '%s' field in '%s'", field, c.connStr)
}
