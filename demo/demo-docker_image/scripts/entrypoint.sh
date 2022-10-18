#!/bin/bash

if [ -z "$NOFLOCK" ]; then
  flock -o "$PGROOT"/mtm-lockfile /bin/sh -c 'NOFLOCK=NOFLOCK /bin/sh /scripts/entrypoint.sh' || exit
  exec sleep infinity
else
  if [ -n "$(ls -A "$PGROOT"/data)" ]; then
    IS_DB_ALREADY_INITIALIZED=1
  fi

  if [ -z "$IS_DB_ALREADY_INITIALIZED" ]; then
    "$PGROOT"/bin/initdb -D "$PGROOT"/data

    echo "log_connections = true
          listen_addresses = '*'
          shared_preload_libraries = 'multimaster'
          wal_level = logical
          max_connections = 1000
          max_prepared_transactions = 3000 # max_connections * N
          max_wal_senders = 100            # at least N
          max_replication_slots = 100      # at least 2N
          wal_sender_timeout = 0
          max_worker_processes = 2000 # (N - 1) * (multimaster.max_workers + 1) + 5" >> "$PGROOT"/data/postgresql.conf

    echo "host replication all all trust
          host mydb all all trust" >> "$PGROOT"/data/pg_hba.conf
  fi

  "$PGROOT"/bin/pg_ctl -D "$PGROOT"/data -l "$PGROOT"/logfile start

  if [ -z "$IS_DB_ALREADY_INITIALIZED" ]; then
    "$PGROOT"/bin/createuser mtmuser -s
    "$PGROOT"/bin/createuser postgres -s
    "$PGROOT"/bin/createdb mydb -O mtmuser
  fi

  exit 0
fi
