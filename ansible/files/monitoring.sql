CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE USER monitoring WITH LOGIN PASSWORD '1234';
ALTER ROLE monitoring SET search_path = mtm, monitoring, pg_catalog, public;
GRANT CONNECT ON DATABASE mydb TO monitoring;
GRANT USAGE ON SCHEMA mtm TO monitoring;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA mtm TO monitoring;
GRANT pg_read_all_settings TO monitoring;
GRANT pg_read_all_stats TO monitoring;
GRANT SELECT ON mtm.cluster_nodes TO monitoring;
