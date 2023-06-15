# Loki Push Pipe

A very simple service to pipe [loki push API](https://grafana.com/docs/loki/latest/api/#push-log-entries-to-loki) message to stdout as JSON object. The objective is to be able to multiplex logs send to loki to another service(s) such as kafka or elasticsearch using already-made utility such as [fluent-bit](https://fluentbit.io/) or [benthos](https://www.benthos.dev/). See [example](./example/docker-compose.yml#L13-L47) for reference on how to pipe received logs to fluent-bit.

## JSON Data

Each entry received by loki-push-pipe will be converted to one json message with the following structure:

```json
{
    "labels": {
        "<key1>":"<value1>",
        "<key2>":"<value2>",
        "...":"..."
    },
    "timestamp": "<rfc3999nano 2006-01-02T15:04:05.999999999Z07:00>",
    "tenant_id": "<the tenant id>",
    "line": "<the actual log>"
}
```

Each entry outputted will be separated by new line.
