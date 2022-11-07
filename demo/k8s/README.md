# Instruction for running k8s for monitoring localy

## Minikube

Install minikube, see [docs](https://minikube.sigs.k8s.io/docs/start/):
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

Run minikube:
```bash
minikube start
```

Add alias for minikube kubectl:
```bash
alias k="minikube kubectl -- "
```

Stop minikube:
```bash
minikube stop
```

Delete minikube:
```bash
minikube delete --all
```

## Install prometheus stack
You need to have _helm_ installed, see [helm installtion](https://helm.sh/docs/intro/install/)

Add helm repo
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

Install prometheus stack:
```bash
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace
```

Delete prometheus:
```bash
helm uninstall -n monitoring prometheus
```

After installing, you can add node port for grafana or start port-forwarding:
```bash
k port-forward -n monitoring svc/prometheus-grafana 30000:80
```
Grafana credentials: admin:prom-operator.

## Start collecting metrics
 
Add monitoring role in your db:
```bash
CREATE USER monitoring WITH LOGIN PASSWORD '1234';
ALTER ROLE monitoring SET search_path = mtm, monitoring, pg_catalog, public;
GRANT CONNECT ON DATABASE demo TO monitoring;
GRANT USAGE ON SCHEMA mtm TO monitoring;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA mtm TO monitoring;
GRANT pg_read_all_settings TO monitoring;
GRANT pg_read_all_stats TO monitoring;
GRANT SELECT ON mtm.cluster_nodes TO monitoring;
```

You need to change **connection URL** for node in [config map](metrics/mtm-metrics-config-cm.yaml)

Create namespace and config that will be mounted to container:
```bash
k create namespace mtm
k apply -f metrics/mtm-metrics-config-cm.yaml -n mtm
```

Deploy mtm-metrics:
```bash
k apply -f metrics/mtm-metrics-deployment.yaml -n mtm
```

Deploy PodMonitor for mtm-metrics if you installed Prometheus stack from above:
```bash
k apply -f metrics/mtm-metrics-pod-monitor.yaml -n mtm
```

**If you edited _config_ and want new changes to take effect:**
```bash
k apply -f metrics/mtm-metrics-config-cm.yaml -n mtm;
k rollout restart deployment mtm-metrics-deployment -n mtm
```

## Grafana dashboard

In Grafana GUI you need to isntall plugin **Statusmap** by _flant_.
Import [dashboards](../../metrics/grafana/)

