package status

import (
	"backup/internal/config"
	"backup/internal/connection"
	"fmt"
	"sort"
)

type MtmStatus struct {
	Cfg         config.Config
	Connections connection.Connections
}

func (m *MtmStatus) Run() {
	statuses := make(map[string]string)
	for _, name := range m.Connections.GetConnNames() {
		db := m.Connections[name].DB
		if db == nil || !db.Ping() {
			statuses[name] = "offline"
			continue
		}

		if !db.HasSharedPreloadLibrary("multimaster") {
			statuses[name] = "multimaster is not in shared_preload_libraries"
			continue
		}

		if !db.HasExtension("multimaster") {
			statuses[name] = "multimaster extension is not installed"
			continue
		}

		if stat := db.MtmStatus(); stat != "" {
			statuses[name] = stat
			continue
		}
	}

	names := make([]string, 0, len(statuses))
	for name := range statuses {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("%s: %s\n", name, statuses[name])
	}

}
