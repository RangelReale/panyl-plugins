# panyl-plugins

Panyl-plugins is a [panyl](https://github.com/RangelReale/panyl) plugin repository.

Panyl is a Golang library to parse logs that may have mixed formats, like log files for multiple services in the same file.

See [panyl-cli-sample](https://github.com/RangelReale/panyl-cli-sample) for a real world usage.

## Plugins

### Metadata

 * `DockerCompose`: extracts the application name from a docker-compose log line
 * `KubeCtlLogs`: extracts the application name from a `kubectl logs --prefix` log line
 * `RubyForeman`: extracts the application name from a ruby foreman log line

### Sequence

 * `RawLine`: treats every line as part of the same sequence

### Parse

 * `ErlangLog`: parses an Erlang log line format
 * `GoLog`: parses a Golang log line format
 * `GoSpaceLog`: parses a Golang `key=value` log line format (`level=info ts=... caller=... msg=...`)
 * `JavaLog`: parses a Java log line format
 * `KubeEventJsonLog`: parses a Kubernetes event JSON format (must be parsed previously by the JSON structure plugin)
 * `MongoLog`: parses a MongoDB log line format
 * `NGINXJsonLog`: parses a NGINX log json format (must be parsed previously by the JSON structure plugin)
 * `NGINXErrorLog`: parses a NGINX error line format
 * `PostgresLog`: parses a Postgres log line format
 * `RedisLog`: parses a Redis log line format
 * `RubyLog`: parses a Ruby log line format

### ParseFormat

  * `ElasticSearchJSON`: parses a ElasticSearch log JSON format

### PostProcess

 * `DebugFormat`: [debugging] shows the format that was detected on each log
 * `DetectJSON`: detects timestamp, level and message on JSON logs that didn't match any known format
 * `Grep`: adds an extra category to logs whose message contains any of the given values

## Author

Rangel Reale (rangelreale@gmail.com)
