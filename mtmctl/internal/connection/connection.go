package connection

import (
	"sort"
)

type Connections map[string]Connection

type Connection struct {
	SSH *SSH
	DB  *DB
}

func (c *Connections) GetConnNames() []string {
	conns := make([]string, 0, len(*c))
	for conn := range *c {
		conns = append(conns, conn)
	}
	sort.Strings(conns)
	return conns
}
