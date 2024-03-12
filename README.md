# Synadia Stats Exporter
At the moment, Synadia cloud is not exporting any prometheus metrics to be consumed
by external tools (like Grafana Cloud). While the Synadia cloud dashboards provide 
a lot of information, we create a per application dashboard in grafana cloud which
shows the progress of consumers and the amount of messages stored in streams.

This exporter will connect to a nats node as a regular client, periodically
consuming stream and consumer information and exporting it as prometheus metrics.

The prometheus metrics are posted on the `/metrics` endpoint at port `2112`

## Usage
```bash
Usage:
  synadia-stats-exporter [flags]

Flags:
      --config string       config file (default is $HOME/.synadia-stats-exporter.yaml)
  -h, --help                help for synadia-stats-exporter
  -i, --interval duration   the interval at which to update the metrics (default 5s)
  -j, --jwt string          the nats user jwt
  -n, --nats string         the nats server url (default "tls://connect.ngs.global")
  -x, --seed string         the nats user seed
```

## Docker
```bash
docker run --rm -it -p 2112:2112 ghcr.io/shono-io/synadia-stats-exporter:main --jwt "<jwt>" --seed "<seed>"
```

## Disclaimer
This is a tool built for a specific problem we endured. While it works for us, it might not work for your
specific use case. We are open to contributions and suggestions, but we do not guarantee any support or
maintenance for this tool.