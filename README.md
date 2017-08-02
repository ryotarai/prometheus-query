# prometheus-query

CLI client to query Prometheus

## Installation

```
$ go get -u github.com/ryotarai/prometheus-query
```

## Usage

```
Usage of prometheus-query:
  -end string
        End time (default to now)
  -format string
        Format (default to json. available formats are json, tsv and csv) (default "json")
  -query string
        Query
  -server string
        Prometheus server URL like 'https://prometheus.example.com' (can be set by PROMETHEUS_SERVER environment variable) (default "")
  -start string
        Start time (default to 1 hour ago)
  -step string
        Step (default to 15s) (default "15s")
```

```
$ PROMETHEUS_SERVER="http://your-prometheus.example.com"
$ QUERY="100 * (1 - avg by(instance_type, availability_zone)(irate(node_cpu{mode='idle'}[5m])))"
```

Output format defaults to JSON:

```
$ prometheus-query -query "$QUERY" | jq .
[
  {
    "time": 1501662831,
    "values": [
      {
        "metric": {
          "availability_zone": "ap-northeast-1b",
          "instance_type": "c3.2xlarge"
        },
        "value": 33.65555555557196
      },
      {
        "metric": {
          "availability_zone": "ap-northeast-1b",
          "instance_type": "c3.4xlarge"
        },
        "value": 32.32500000000012
      },
...
```

CSV and TSV formatters are available:

```
$ prometheus-query -query "$QUERY" -format csv
time,availability_zone:ap-northeast-1b|instance_type:c3.2xlarge,availability_zone:ap-northeast-1b|instance_type:c3.4xlarge
1501662920,36.910185,26.927778
1501662935,34.331481,27.270139
...
```

Query range can be specified by `-start` and `-end` options and they are parsed by https://github.com/ymotongpoo/datemaki .
Start time defaults to 1 hour ago and end time defaults to now:

```
$ prometheus-query -query "$QUERY" -start-time '2 day ago' -end-time '1 day ago'
```
