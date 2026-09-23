# Code Review — panyl-plugins

Scope: every non-test Go file in `metadata/`, `sequence/`, `parse/`, `parseformat/` and `postprocess/` at `master` (`f62c3a6`), checked against how panyl v2.1.4 calls plugins (`job.go`). `go build`, `go vet` and `go test ./...` are clean.

I confirmed every bug in the **Bugs** section by running it: a throwaway test called the plugin with the input shown and printed the result. That test was deleted afterwards.

## Summary

| # | Severity | Location | Issue |
|---|----------|----------|-------|
| 1 | High | `parse/javalog.go:20` | Regex accepts any word as the level, so it takes lines meant for other parsers (e.g. Postgres) |
| 2 | High | `parse/golog.go:25` | Greedy caller group eats part of the message when the message contains `something.go:N` |
| 3 | Medium | `parse/nginx_json_log.go:70` | Wrong variable in the condition, so a trailing ` -- ` is always appended |
| 4 | Medium | `metadata/docker_compose.go:54` | Always drops the first character after `\|` |
| 5 | Medium | `parse/erlanglog.go:61-69` | `critical`/`alert`/`emergency` are reported as INFO |
| 6 | Medium | `postprocess/grep.go:42` | An empty or whitespace-only value matches every line |
| 7 | Medium | several parsers | Timestamps without a zone are read as UTC, but those apps usually log local time |
| 8 | Low | `parse/gospacelog.go:70-82` | Level matching is case-sensitive, so no level is set for uppercase or unknown levels |
| 9 | Low | `metadata/kubectl_logs.go:92-115` | The pod-name heuristic can return an empty name or cut a real name segment |
| 10 | Low | `metadata/ruby_foreman.go:40` | `TrimSpace` removes indentation (e.g. from stack traces) |
| 11 | Low | `parse/nginx_json_log.go:12` | Format name is `cb_go_json_log` (copy-paste) |

The **Cleanups** section below lists smaller items.

---

## Bugs

### 1. `JavaLog` takes lines meant for other parsers — `parse/javalog.go:20`

```go
var javaLogRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} [^\s]+)\s+(\w+)\s+(.*)$`)
```

Group 2 accepts any word. panyl stops at the first parse plugin that matches, in the order plugins were registered (`job.go:143-154`). So if `JavaLog` is registered before `PostgresLog`, it takes every Postgres line:

```
input:  2022-04-05 14:29:07.500 UTC [73] ERROR:  relation "users" does not exist
result: matched=true, Data["level"]="UTC", MetadataLevel=""   (ERROR is lost)
```

Other `YYYY-MM-DD hh:mm:ss …` formats are hit the same way. **Fix:** only accept known levels, e.g. `(TRACE|DEBUG|INFO|WARN|WARNING|ERROR|FATAL|SEVERE)`. This also adds `TRACE`/`FATAL` to the level mapping (lines 58-66), which currently leaves them unset.

### 2. `GoLog` caller group is greedy — `parse/golog.go:25`

```go
goLogRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T[^\s]+)\s+(\w+)\s+(.*\.go:\d+)\s+(.*)$`)
```

`(.*\.go:\d+)` runs to the **last** `.go:N` on the line:

```
input:   2022-03-10T19:53:21.434Z	INFO	app/main.go:10	failed loading config.go:42 retrying
source:  "app/main.go:10\tfailed loading config.go:42"
message: "retrying"
```

When `SourceAsCategory` is on, the category is wrong too. **Fix:** `(\S+\.go:\d+)`. The unused `GoSpaceLogRe` (`parse/gospacelog.go:27`) is a copy of this regex; see the Cleanups section.

### 3. `NGINXJsonLog` checks the wrong variable — `parse/nginx_json_log.go:70`

```go
if logmessage := item.Data.StringValue("message"); message != "" {
    message = fmt.Sprintf("%s -- %s", message, logmessage)
}
```

The condition checks `message`, which is never empty at this point, instead of `logmessage`. Every log without a `message` field ends with a dangling separator:

```
result: "GET / [status:200] -- "
```

**Fix:** `if logmessage := …; logmessage != "" {`.

### 4. `DockerCompose` drops a character — `metadata/docker_compose.go:53-54`

```go
if len(item.Line) > matches[1] {
    item.Line = item.Line[matches[1]+1:]
```

`matches[1]` is already the index just past `|`, and the code skips one more byte on the assumption that it is a space. If no space follows, real content is lost:

```
input:  "app |message"
result: "essage"
```

**Fix:** slice at `matches[1]` and remove at most one leading space, e.g. `strings.TrimPrefix(item.Line[matches[1]:], " ")`. Better still, put the optional space in the regex (`…\| ?`).

### 5. `ErlangLog` reports critical levels as INFO — `parse/erlanglog.go:61-69`

Erlang `logger` levels are `emergency, alert, critical, error, warning, notice, info, debug`. Only `debug`, `warning` and `error` are mapped, and the `else` branch sends everything else to INFO:

```
input:  2025-09-04T13:26:52.315705+00:00 [critical] boom
result: MetadataLevel="info"
```

**Fix:** map `emergency|alert|critical` to ERROR.

### 6. `Grep` with an empty value matches every line — `postprocess/grep.go:42`

```go
strings.Contains(strings.ToLower(message), strings.ToLower(strings.TrimSpace(value)))
```

A value of `""` or `" "` trims to `""`, and `strings.Contains(x, "")` is always true. So every line gets tagged:

```
Values: ["foo", " "]  →  ExtraCategories: ["[grep][ ]"]
```

This is easy to hit from CLI input like `--grep "a,,b"`. **Fix:** skip values that are empty after trimming. Also see the efficiency notes on this function in the Cleanups section.

### 7. Timestamps without a zone are read as UTC

`time.Parse` with a layout that has no zone returns UTC. These formats have no zone and are normally written in the server's **local** time:

- `RubyLog` (`rubylog.go:23`)
- `JavaLog` (`javalog.go:22`)
- `RedisLog` (`redislog.go:23`)
- `NGINXErrorLog` (`nginxerrorlog.go:24`)

On a non-UTC host, the times are off by the UTC offset. For example, `sample/sample.go` calls `.Local()` on these UTC values, which shifts them a second time. This is a design choice rather than a local bug. **Suggestion:** add an optional `Location *time.Location` field to these plugins (defaulting to the current behavior) and use `time.ParseInLocation`.

Two related points:
- `PostgresLog` hard-codes ` UTC ` in its regex (`postgres_log.go:22`), so servers with any other `log_timezone` are not matched at all.
- `postgresTimestampFormat` requires exactly three decimal digits.

### 8. `GoSpaceLog` level matching is case-sensitive — `parse/gospacelog.go:70-82`

```
input:  level=WARN ts=2024-12-18T14:55:27Z msg="x"
result: MetadataLevel=""
```

logfmt producers vary in case (`WARN`, `Info`), and `fatal`/`panic`/`crit` are not mapped at all. **Fix:** compare against `strings.ToLower(level)` and add the missing levels, as `postprocess/detect_json.go:49` already does.

### 9. `KubeCtlLogs.parsePodName` edge cases — `metadata/kubectl_logs.go:92-115`

The heuristic assumes Deployment-style pod names (`<name>-<rs-hash>-<suffix>`). Stripping hex-looking hashes is correct, because Kubernetes' pod-template-hash maps decimal digits into `[4-9bcdf]`. Two other assumptions are not:

- `"deadbeef-abcd"` → strips `abcd`, then strips `deadbeef` → `ns` is empty → returns `"/c"`.
- Any last segment of 4-5 characters is treated as the random suffix, so a plain pod `api-proxy` becomes `"api/c"`.

**Suggestion:** never strip the last remaining segment. Only strip a 4-5 character suffix when it matches the Kubernetes random alphabet (`^[bcdfghjklmnpqrstvwxz2456789]{5}$`). Kubernetes suffixes are exactly 5 characters, so 4-character segments should not be stripped either.

### 10. `RubyForeman` removes indentation — `metadata/ruby_foreman.go:40`

```go
text := strings.TrimSpace(matches[3])
```

```
input:  "16:41:59 api.1   |     at foo.rb:1"
result: "at foo.rb:1"
```

Stack-trace indentation is lost, and later sequence or parse plugins can no longer see it. **Fix:** remove only the single separator space (`strings.TrimPrefix(matches[3], " ")`).

### 11. `NGINXJsonLog` format name — `parse/nginx_json_log.go:12`

```go
const NGINXJsonLogFormat = "cb_go_json_log"
```

This looks copied from another plugin. Consumers that switch on `MetadataFormat` see a "go json" format for NGINX lines. `KubeEventJsonLogFormat = "cb_kube_json_log"` has the same unexplained `cb_` prefix. Renaming either one breaks consumers that match on these strings, so change it deliberately, e.g. in a new version.

---

## Cleanups

- **Dead code in `GoSpaceLog`:** the exported `GoSpaceLogRe` (`gospacelog.go:27`) and `encodeFieldsExcept` (`gospacelog.go:147`) are never used. `GoSpaceLogRe` is also part of the public API, so removing it later is a breaking change.
- **Misnamed variable:** `backTick` in `GoSpaceLog.splitString` (`gospacelog.go:134`) actually tracks a backslash escape.
- **Unused field:** `RawLine.SourceAsCategory` (`sequence/rawline.go:11`) is never read. The doc comment "returns all lines" also doesn't say that `BlockSequence` always returning `true` means each line becomes its own item.
- **`Grep` efficiency:** `PostProcess` (`postprocess/grep.go:25-58`) has three problems:
  - It rebuilds `message`, including a full `json.Marshal(item.Data)`, once per value. Build it once before the loop.
  - It lowercases the message once per value.
  - It calls `item.Clone()` (a deep copy) on every match just to read the category, and the clone is then discarded. Read `item.Metadata` directly.
- **`Grep` assertion:** `var _ panyl.PluginPostProcess = (*Grep)(nil)` checks the pointer type, but all methods have value receivers and every other plugin asserts the value type. Use `Grep{}`.
- **Duplicated whitelist check:** `DockerCompose`, `KubeCtlLogs` and `RubyForeman` each have the same whitelist loop. `slices.Contains(m.ApplicationWhitelist, application)` replaces each one (`kubectl_logs.go` already imports `slices`).
- **Dead branch in `KubeCtlLogs`:** at `kubectl_logs.go:77`, `len(item.Line) > len(matches[1])` is always true, because the line also contains `[`, `]` and a space. The `else` branch never runs.
- **Redundant parse in `DetectJSON`:** the `time.RFC3339Nano` fallback (`detect_json.go:36`) is never reached, because `time.Parse(time.RFC3339, …)` already accepts fractional seconds. Numeric epoch timestamps (`"ts": 1690000000.1`) are not handled.
- **Same post-process order:** `DetectJSON` and `Grep` both use `PostProcessOrderFirst + 1`. Their relative order then depends on registration order (`job.go:394`). If `Grep` runs first on a JSON line, it searches the marshalled JSON instead of the detected message. Give `Grep` a later order if the detected message is the intended search target.
- **Unescaped dot:** `RedisLogRe` has `.\d{3}` (`redislog.go:22`), which matches any character. Use `\.\d{3}`.
- **Unanchored regex:** `rubyLogRe` has no `^` (`rubylog.go:22`), so it can match in the middle of a line. If that is intentional (e.g. to tolerate prefixes), say so in a comment.
- **Missing levels:**
  - `GoLog` doesn't map zap's `DPANIC`/`PANIC`/`FATAL`.
  - `MongoLog` doesn't map `F` (fatal).
  - MongoDB 4.4+ writes structured JSON logs, which `MongoLog` doesn't handle.
- **Strict layout in `KubeEventJsonLog`:** `kubeEventTimestampFormat = "2006-01-02T15:04:05Z"` (`kube_event_json_log.go:26`) only accepts a literal `Z`. `time.RFC3339` accepts that and offsets too.
- **Plugin in the wrong package:** `KubeEventJsonLog` and `NGINXJsonLog` implement `PluginParseFormat` but live in `parse/` rather than `parseformat/`. Moving them is a breaking change, so at least note it where they are defined.
- **Doc typos:**
  - `PostgresLog` comment says "parses MongoDB log lines format" (`postgres_log.go:13`).
  - `GoSpaceLog` comment says "Golang log lines" but the plugin parses logfmt.
  - `RubyForeman` comment says "roby foreman".
- **Tests:** there are no tests for `parseformat/`, `postprocess/`, `sequence/` or `metadata/kubectl_logs.go`. Each bug above has a one-line input that would make a good table case.
