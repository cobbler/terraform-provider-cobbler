## 6.0.0

IMPROVEMENTS

* Targets Cobbler 4.0.x via `cobblerclient` v1.0.0. Provider startup now refuses
  Cobbler < 4.0.0 with a `cobbler server too old` diagnostic.
* New `cobbler_template` resource and data source backed by the first-class
  Cobbler 4.0.0 `Template` item type (replaces `cobbler_snippet` and
  `cobbler_template_file`).
* New `cobbler_network_interface` resource and data source for per-interface
  CRUD; supports nested `ipv4`, `ipv6`, and `dns` option objects matching
  upstream `IPv4Option` / `IPv6Option` / `DNSInterfaceOption`.
* New `cobbler_distro_group`, `cobbler_profile_group`, `cobbler_system_group`
  resources and data sources for bulk-operation patterns.
* New `cobbler_system.uid` Computed attribute exposing the server-assigned UID;
  required to wire `cobbler_network_interface.system`.

BACKWARDS INCOMPATIBILITIES

* `cobbler_snippet` and `cobbler_template_file` removed — migrate to
  `cobbler_template` (see `MIGRATION.md`).
* `cobbler_system.interface = { ... }` map removed — split each interface into
  its own `cobbler_network_interface` resource using the `<ifname>@<systemname>`
  name syntax.
* `mgmt_classes` and `mgmt_parameters` attributes removed from `cobbler_distro`,
  `cobbler_image`, `cobbler_menu`, `cobbler_profile`, `cobbler_system` (the
  `MgmtClass` item type was removed upstream in 4.0.0).
* Minimum Cobbler server: 4.0.0. Users on 3.3.x must stay on v5.x.

## 3.0.0 (Jan 27, 2022)

IMPROVEMENTS

* Supports Cobbler: v3.3.x
* Moved test harness to local docker container for easier\faster development

BACKWARDS INCOMPATIBILITIES

* Rewrites to support Cobbler 3.3.x (will break support for Cobbler 3.2.x and older (EOL)).
* `next_server` attribute is now either `next_server_v4` or `next_server_v6`
* `boot_loader` string attribute is renamed to `boot_loaders` and changed from a string to a list
* The following string attributes are now lists: `fetchable_files`, `kernel_options`, `kernel_options_post`, 
`template_files`, `autoinstall_meta`, and `repos`

## 2.0.1 (April 30, 2020)

BUG FIXES

* Bugfix in dependency "cobblerclient" - IPv6DefaultGateway

## 2.0.0 (March 02, 2020)

BACKWARDS INCOMPATIBILITIES

* Rewrites to support Cobbler 3.x (will break support for Cobbler 2.x (EOL)).

## 1.1.1 (Unreleased)

## 1.1.0 (June 07, 2019)

IMPROVEMENTS

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 1.0.1 (February 22, 2018)

FEATURES:

* Support for self-signed certificates ([#11](https://github.com/terraform-providers/terraform-provider-cobbler/issues/11))

BUG FIXES

* Recreate systems if they were deleted outside of Terraform ([#14](https://github.com/terraform-providers/terraform-provider-cobbler/issues/14))

## 1.0.0 (November 15, 2017)

FEATURES:

__New Resource:__ `cobbler_repo` ([#3](https://github.com/terraform-providers/terraform-provider-cobbler/issues/3))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
