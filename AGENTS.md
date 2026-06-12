# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Terraform / OpenTofu provider for [Cobbler](https://cobbler.github.io/). Built on
`terraform-plugin-framework` (protocol v6) — **not** `terraform-plugin-sdk/v2`. Wraps the
`github.com/cobbler/cobblerclient` XML-RPC client. Requires a running Cobbler server (≥ 3.3.0) to do
anything useful; the provider's `Configure` calls `client.Login()` at provider start, so even unit-like flows
that go through `Configure` need a reachable server.

Module path: `github.com/cobbler/terraform-provider-cobbler`. Entry point is `main.go`, served at
`registry.terraform.io/cobbler/cobbler`. All provider code lives under `internal/` (consumers cannot import it).

## Commands

```bash
make build       # go install — formats first via fmtcheck
make fmt         # gofmt -w on all non-vendor .go files
make fmtcheck    # CI gate; fails if anything is unformatted
make vet         # go vet
make errcheck    # unchecked-error scan via kisielk/errcheck
make test        # unit tests, 30s timeout, parallel=4
make testacc     # full acceptance suite — boots Cobbler in Docker first (see below)
make docs        # regenerate docs/ via tfplugindocs (must be installed)
```

Run a single test:

```bash
# Unit test in one package
go test -run TestStringFrom_Inherited ./internal/inherit/

# Single acceptance test (env vars must be set; Cobbler must already be running)
TF_ACC=1 TF_ACC_PROVIDER_NAMESPACE=cobbler \
  COBBLER_URL=http://localhost:8081/cobbler_api \
  COBBLER_USERNAME=cobbler COBBLER_PASSWORD=cobbler \
  go test -v -run TestAccCobblerDistro_basic ./internal/distro/
```

`TF_ACC_PROVIDER_NAMESPACE=cobbler` is required — without it the framework uses the legacy `-` wildcard
namespace which OpenTofu rejects.

### Acceptance test environment

`make testacc` shells out to `docker/start.sh`, which:

1. Clones a pinned Cobbler commit into `docker/cobbler_source/` (build context for the image).
2. Downloads the Ubuntu 20.04 legacy-server ISO (~860 MB, cached at repo root) and extracts it with
   `xorriso` into `extracted_iso_image/` (used by the import-distro tests).
3. Builds the `cobbler-dev` image and brings up `docker/compose.yml`, exposing the API on
   `http://localhost:8081/cobbler_api`.

`xorriso` must be installed on the host (`apt-get install -y xorriso` etc.). The ISO and extracted dir
are recreated only when their checksums no longer match — safe to leave in place between runs. CI tests
across a matrix of Cobbler 3.3.x commits and both Terraform and OpenTofu (`.github/workflows/testing.yml`).

## Architecture

### Per-resource package layout

Every Cobbler object type is its own package under `internal/` with the same five-file shape:

```text
internal/<thing>/
    resource.go              # schema + Create/Read/Update/Delete/ImportState
    resource_model.go        # tfsdk-tagged struct + conversion helpers to/from cobblerclient
    resource_test.go         # TestAccCobbler<Thing>_* acceptance tests
    data_source.go           # read-only data source
    data_source_model.go
    data_source_test.go
```

Current set: `distro`, `image`, `menu`, `profile`, `repo`, `snippet`, `system`, `template_file`. To add a
new resource, copy an existing package, register `NewResource` / `NewDataSource` in
`internal/provider/provider.go` (both the `Resources` and `DataSources` lists), and add docs under
`docs/resources/` and `docs/data-sources/`.

The provider's `Configure` stores a `*client.Config` (with the logged-in `cobbler.Client` attached) into
both `resp.ResourceData` and `resp.DataSourceData`. Each resource's `Configure` type-asserts to
`*clientpkg.Config` and keeps the embedded `cobbler.Client`.

### Inheritable fields — `internal/inherit/`

Cobbler distinguishes "value not set, inherit from parent" from "explicitly empty". In v4 the schema
represents every inheritable field as a **nested object** with `{value, inherited}` (see `MIGRATION.md` for
the v3→v4 attribute shape change). The `inherit` package provides `From`/`To` helpers per primitive type
(`bool`, `int`, `float64`, `string`, `stringlist`, `stringmap`) that convert between
`cobbler.Value[T]{Data, IsInherited}` and the Terraform `types.Object`. New inheritable fields should reuse
these helpers rather than hand-rolling the conversion.

When `inherited = true`, write the field's null/zero value into state for `value`. The Optional+Computed
fields also need the Unknown-in-plan handling described under **Provider-framework gotchas** below.

### System resource specifics

- `internal/system/interface.go` — the `interface` attribute is a **map** keyed by interface name (not the
  v3 set-of-blocks). The map key replaces the old `name` sub-attribute.
- `internal/system/mutex.go` — a package-level `sync.Mutex` (`systemSyncLock`) serializes Cobbler system
  mutations because the upstream API is not safe for concurrent edits to the same object.

### Test helpers — `internal/acctest/`

`acctest.ProtoV6ProviderFactories` is the map every `resource.TestCase` should use. `acctest.PreCheck(t)`
verifies `COBBLER_URL` / `COBBLER_USERNAME` / `COBBLER_PASSWORD` are set. The package's `init()` also
constructs a bare `CobblerApiClient` for tests that need to assert side effects directly against the API
(e.g. CheckDestroy).

## Conventions

- Go ≥ 1.25 (see `go.mod`); minimum Terraform 1.0 / OpenTofu 1.6.
- `make fmtcheck` runs in `build` — keep `gofmt -l` clean or CI fails.
- Lint config in `.golangci.yml` (golangci-lint v2). Don't lint files under `go/` (excluded).
- Markdown line length is 120 (`.markdownlint.yml`); `docs/` is generated by `tfplugindocs` so edit
  `templates/` (or the resource schema `Description` strings) rather than the rendered files.
- Docs in `docs/` are generated — don't hand-edit; regenerate with `make docs`.
- Schema `Description` strings end up in the rendered docs verbatim — write them as user-facing prose.
- The vendored `extracted_iso_image/` and the `.iso` at repo root are test fixtures; never commit changes to
  them and never delete them as "cleanup" (they're cached intentionally).

## Provider-framework gotchas

### Optional+Computed attributes are Unknown in the plan, not the prior state value

In `terraform-plugin-framework` v1.19.0, an `Optional+Computed` attribute that is **not set in config**
arrives in the plan as **Unknown** — not null, and not the previously-stored state value — for BOTH Create
and Update. The framework log line is `"marking computed attribute that is null in the config as unknown"`.

Practical consequences when transforming the plan into a `cobblerclient` call:

1. **Enum strings (e.g. `virt_disk_driver`, `virt_type`, `interface_type`)** — `types.String.ValueString()`
   returns `""` for an Unknown value. Cobbler's enum validation rejects `""`. Use the existing
   `stringOrInherit` helper (see `internal/image/resource.go`, `internal/profile/resource.go`) which maps
   empty/null/unknown to the sentinel `"<<inherit>>"` that Cobbler understands:

   ```go
   func stringOrInherit(s types.String) string {
       if v := s.ValueString(); v != "" {
           return v
       }
       return "<<inherit>>"
   }
   ```

2. **List / map fields** — guard `ElementsAs` with `!IsNull() && !IsUnknown()` or you get a panic
   `received unknown value, target type cannot handle unknown values`:

   ```go
   var list []string
   if !data.SomeList.IsNull() && !data.SomeList.IsUnknown() {
       diags.Append(data.SomeList.ElementsAs(ctx, &list, false)...)
   }
   ```

3. **Bools** — `.ValueBool()` returns `false` for Unknown without warning. On an import followed by apply,
   that silently clears flags such as `management` and `static` on system interfaces. Always check
   `IsNull() || IsUnknown()` before reading.

4. The framework **does not** preserve the prior state value for Optional+Computed attributes during
   planning. Add `UseStateForUnknown()` plan modifiers on every Optional+Computed attribute (string, bool,
   list, map, single-nested-object) — otherwise every plan shows them as `(known after apply)` even when
   nothing changed.

**Cobbler-client asymmetry:** `cobblerclient.CreateProfile` / `CreateSystem` have `if field == "" { field = inherit }`
fallbacks, but the corresponding `UpdateProfile` / `UpdateSystem` calls do not. The provider must do the
empty→inherit translation for enum fields itself on the update path.
