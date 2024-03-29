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
