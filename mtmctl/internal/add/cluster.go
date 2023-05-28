package add

import (
	"fmt"
	"mtmctl/internal/config"
	"mtmctl/internal/connection"
	"strconv"
	"strings"
	"sync"
)

type MtmAddCluster struct {
	Cfg            config.Config
	SSHConnections connection.Connections
	Nodes          []string
}

func (m *MtmAddCluster) Run() {
	m.reinitData()

	m.prepareNodes()

	m.createCluster()

	m.verifyCluster()

	m.addMonitoring()

	fmt.Println("Cluster created")
}

func (m *MtmAddCluster) reinitData() {
	var wg sync.WaitGroup

	for _, nodename := range m.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()

			nodeSSH := m.SSHConnections[node].SSH

			nodeSSH.PgCtlStop()
			nodeSSH.RemovePGDATA()

			nodeSSH.Initdb()

			fmt.Printf("Reinitialized database on '%s'\n", node)
		}(nodename)
	}

	wg.Wait()
}

func (m *MtmAddCluster) prepareNodes() {
	pgsql_settings := []string{
		`"listen_addresses = '*'"`,
		`"shared_preload_libraries = 'multimaster'"`,
		`"wal_level = logical"`,
		`"max_connections = 100"`,
		`"max_prepared_transactions = 300"`,
		`"max_wal_senders = 10"`,
		`"max_replication_slots = 10"`,
		`"wal_sender_timeout = 0"`,
		`"max_worker_processes = 500"`,
	}

	pg_hba_settings := []string{
		`"host replication mtmuser 0.0.0.0/0 md5"`,
		`"host all all 0.0.0.0/0 md5"`,
	}

	var wg sync.WaitGroup

	for _, nodename := range m.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()

			fmt.Printf("Preparing node '%s'\n", node)

			nodeSSH := m.SSHConnections[node].SSH

			for _, line := range pgsql_settings {
				nodeSSH.AddPostgresqlConfLine(line)
			}

			for _, line := range pg_hba_settings {
				nodeSSH.AddPgHBAConfLine(line)
			}

			nodeSSH.PgCtlStart()
			nodeSSH.RunPSQL("CREATE USER mtmuser WITH SUPERUSER PASSWORD '1234'")
			nodeSSH.RunPSQL("CREATE DATABASE mydb OWNER mtmuser")
		}(nodename)
	}

	wg.Wait()

	fmt.Println("All nodes are prepared")
}

func (m *MtmAddCluster) createCluster() {
	nodeSSH := m.SSHConnections["node1"].SSH

	fmt.Println("Creating cluster")
	nodeSSH.RunPSQL_MTM("CREATE EXTENSION multimaster")
	nodeSSH.RunPSQL_MTM(`SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=node1', '{\"dbname=mydb user=mtmuser host=node2\", \"dbname=mydb user=mtmuser host=node3\"}');`)

	fmt.Print("Waiting for cluster to become available")
	nodeSSH.WaitForCluster(len(m.Nodes))
}

func (m *MtmAddCluster) verifyCluster() {
	nodeSSH := m.SSHConnections.GetFirstConnection().SSH

	fmt.Print("Checking nodes availability: ")

	out := nodeSSH.RunPSQL_MTM("SELECT count(*) FROM mtm.nodes() WHERE enabled = 't' and connected = 't'")
	count, err := strconv.Atoi(out)
	if err != nil {
		fmt.Printf("error: %s \n", out)
	}

	if count == len(m.Nodes) {
		fmt.Printf("all %d nodes are enabled and connected\n", len(m.Nodes))
	} else {
		fmt.Printf("only %s out of %d nodes are enabled and connected\n", out, len(m.Nodes))
	}
}

func (m *MtmAddCluster) addMonitoring() {
	nodeSSH := m.SSHConnections.GetFirstConnection().SSH

	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS pg_stat_statements`,
		`CREATE USER monitoring WITH LOGIN PASSWORD '1234'`,
		`ALTER ROLE monitoring SET search_path = mtm, monitoring, pg_catalog, public`,
		`GRANT CONNECT ON DATABASE mydb TO monitoring`,
		`GRANT USAGE ON SCHEMA mtm TO monitoring`,
		`GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA mtm TO monitoring`,
		`GRANT pg_read_all_settings TO monitoring`,
		`GRANT pg_read_all_stats TO monitoring`,
		`GRANT SELECT ON mtm.cluster_nodes TO monitoring`,
	}

	nodeSSH.RunPSQL_MTM(strings.Join(queries, "; "))
}
