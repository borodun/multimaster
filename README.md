# multimaster
Postgres multimaster \
Deploy modified postgres and configure everything
```shell
ansible-playbook site.yml -i hosts  --ask-become-pass -t prepare,build,install
```
_BECOME password_ is password for superuser on hosts

Start/stop/status postgres
```shel
ansible-playbook site.yml -i hosts -t start/stop/status
```

Other tags
```shel
ansible-playbook site.yml -i hosts --list-tags
```

On one node start multimaster cluster
```shell
psql -U mtmuser -d mydb
```
```sql
CREATE EXTENSION multimaster;
SELECT mtm.init_cluster('dbname=mydb user=mtmuser host=localhost',
    '{"dbname=mydb user=mtmuser host=node1", "dbname=mydb user=mtmuser host=node2"}');
```
Check
```sql
SELECT * FROM mtm.status();
SELECT * FROM mtm.nodes();
```
