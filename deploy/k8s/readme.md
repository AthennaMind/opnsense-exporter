# Simple Kubernetes Deployment

Included are two files:

- `deployment.yaml`: sets up a bare-bones deployment for just the exporter and a service to expose it to the rest of the cluster.
- `scrape.yaml`: a [`ScrapeConfig`](https://prometheus-operator.dev/docs/user-guides/scrapeconfig/) CRD which will configure Prometheus to scrape metrics from the exporter

Both files are well commented and should be easy to modify to suit your needs.

## Pre-requisites

When you [generate API keys](https://docs.opnsense.org/development/how-tos/api.html#creating-keys) for an OPNSense user, you will get a `.txt` file with the API key and secret.

> **Note**
> Consider setting up a proper group with limited API permissions.
> If you use the `root` user to generate the API key, the key will have full access to the OPNSense API.

We'll add a few more configuration directives to this file and use it to create a secret in your Kubernetes cluster.

To start with, the file containing the key/secret should look like this:

```shell
❯ cat opnsense_apikey.txt
key=xt<...>Nt
secret=EK<...>ho
```

Add both `host`, `protocol` directives to the file:

```shell
❯ cat opnsense_apikey.txt
key=xt<...>Nt
secret=EK<...>ho
# Your OPNSense host name/IP here
host=opnsense.lan
protocol=https
```

Then create the secret in your Kubernetes cluster:

```shell
❯ kubectl create secret generic opnsense-exporter-cfg --namespace=o11y --from-env-file=opnsense_apikey.txt 
secret/opnsense-exporter-cfg created
```

With the secret created, the exporter can be deployed.

```shell
❯ k apply -f deployment.yaml
deployment.apps/opnsense-exporter created
service/opnsense-exporter unchanged
```

Check your work:

```shell
❯ kubectl run debug --rm -i --tty --restart=Never --image=alpine --namespace=o11y
<...>
/ # wget --quiet -O- opnsense-exporter.o11y.svc.cluster.local:8080/metrics
<...>
```
