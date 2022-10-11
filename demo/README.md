# Demo instruction

1. Start multimaster cluster, see [guide](../local-deploy/README.md)
2. Build and run _mtm-connector_, see [guide](../mtm-connector/README.md)
3. Build and copy to phone _mtm-joiner_, see [guide](../mtm-joiner/README.md)
4. Add nodeon phone:
```bash
./mtm-joiner -u <connector-url>
```
5. Drop node on phone:
```bash
./mtm-joiner -u <connector-url> --drop
```
6. On PC clean up after dropping, see [guide](../local-deploy/README.md):
```bash
./scripts/clean_up.sh
```
