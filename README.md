# Loki Push Pipe

A very simple service to pipe [loki push API](https://grafana.com/docs/loki/latest/api/#push-log-entries-to-loki) message to stdout as JSON object. The objective is to be able to multiplex logs send to loki to another service(s) such as kafka or elasticsearch using already-made utility such as [fluent-bit](https://fluentbit.io/) or [benthos](https://www.benthos.dev/). See [example](./example/docker-compose.yml#L13-L47) for reference on how to pipe received logs to fluent-bit.
