# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A library of plugins for [panyl](https://github.com/RangelReale/panyl) (`github.com/RangelReale/panyl/v2`), a Go library that parses logs with mixed formats (e.g. several services' logs interleaved in one stream). This repo has no binary: it is imported as `github.com/RangelReale/panyl-plugins/v2`. `sample/sample.go` is a runnable example that wires plugins into a `panyl.Processor` and reads from stdin; [panyl-cli-sample](https://github.com/RangelReale/panyl-cli-sample) is the real-world consumer.

## Commands

```sh
go build ./...                          # build (CI runs go build -v ./... and go test -v ./... on Go 1.23)
go test ./...                           # all tests (also: task test)
go test ./parse -run TestGoLog -v       # a single test
go run ./sample < some.log              # try plugins against a log file
```

Releases are git tags: `task release-version VERSION=vX.Y.Z` (needs a clean tree) tags and pushes, which triggers GoReleaser in `.github/workflows/release.yml`.

## Architecture

Packages are grouped by the panyl plugin interface they mainly implement. The panyl processor calls each kind at a different stage of the pipeline:

- `metadata/` — `panyl.PluginMetadata` (`ExtractMetadata`): runs on each raw line, pulls prefix info such as the application name (docker-compose `app |`, foreman, kubectl), and strips that prefix from `item.Line`. These plugins usually also implement `panyl.PluginSequence` (`BlockSequence`) so that a change of application ends a multi-line block.
- `sequence/` — `panyl.PluginSequence` only: decides whether consecutive lines belong to the same item.
- `parse/` — mostly `panyl.PluginParse` (`ExtractParse(ctx, lines panyl.ItemLines, item)`): matches raw text lines, usually with a package-level `regexp`. Most accept only `len(lines) == 1`. On a match, the plugin calls `item.MergeLinesData(lines)`, clears `item.Line`, fills `item.Data`, and sets `item.Metadata`. Some files here (`nginx_json_log.go`, `kube_event_json_log.go`) instead implement `panyl.PluginParseFormat`: see below.
- `parseformat/` — `panyl.PluginParseFormat` (`ParseFormat(ctx, item)`): runs after a structure plugin (e.g. panyl's `structure.JSON{}`) has already put data in `item.Data`. It checks `item.Metadata[panyl.MetadataStructure] == panyl.MetadataStructureJSON` plus the expected keys, then sets the metadata.
- `postprocess/` — `panyl.PluginPostProcess` (`PostProcess`, `PostProcessOrder`): runs last on fully parsed items (debugging, grep filtering, JSON detection).

What a plugin must provide:
- An `IsPanylPlugin()` marker method.
- A compile-time interface assertion, e.g. `var _ panyl.PluginParse = GoLog{}`.
- An exported `<Name>Format` string constant, written to `item.Metadata[panyl.MetadataFormat]`.
- Standard metadata keys where it can fill them: `panyl.MetadataTimestamp` (as `time.Time`), `MetadataLevel` (mapped to the `panyl.MetadataLevel*` constants), `MetadataMessage`, `MetadataCategory`, `MetadataApplication`.

Plugins are value-type structs, and their option fields (e.g. `DisableCaller`, `SourceAsCategory`) change how they match or what they emit.

Tests sit next to each plugin as `*_test.go`. They are table-driven with testify `assert`: build an item with `panyl.InitItem()`, call the plugin method directly (e.g. `ExtractParse(ctx, panyl.ItemLines{&panyl.Item{Line: src}}, item)`), then assert on `item.Metadata`.

When adding a plugin, also list it in the README's "Plugins" section.
