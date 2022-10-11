# Info
Script for deploying multimaster localy
## Usage
You need to change node adresses in _resources/init_mm1.sql_

```bash
# Start all nodes
./sctipts/full_bootstart_running.sh
# Stop, reinstall and start cluster
./restart_all.sh
#  Stop/stop/status all nodes
./scripts/poke_all.sh start/stop/status
# Clean up cluster after dropping node
./scripts/clean_up.sh
```
