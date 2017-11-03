---
layout: "cobbler"
page_title: "Cobbler: cobbler_repo"
sidebar_current: "docs-cobbler-resource-repo"
description: |-
  Manages a repo within Cobbler.
---

# cobbler_repo

Manages a repo within Cobbler.

## Example Usage

```hcl
resource "cobbler_repo" "my_repo" {
  name           = "my_repo"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["trusty"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}
```

## Argument Reference

The following arguments are supported:

* `apt_components` - (Optional) List of Apt components such as main,
  restricted, universe. Applicable to apt breeds only.

* `apt_dists` - (Optional) List of Apt distribution names such as trusty,
  trusty-updates. Applicable to apt breeds only.

* `arch` - (Optional) The architecture of the repo. Valid options
  are: i386, x86_64, ia64, ppc, ppc64, s390, arm.

* `breed` - (Required) The "breed" of distribution. Valid options
  are: rsync, rhn, yum, apt, and wget. These choices may vary depending on the
  version of Cobbler in use.

* `comment` - (Optional) Free form text description.

* `createrepo_flags` - (Optional) Flags to use with `createrepo`.

* `environment` - (Optional) Environment variables to use during repo command
  execution.

* `keep_updated` - (Optional) Update the repo upon Cobbler sync. Valid values
  are true or false.

* `mirror` - (Required) Address of the repo to mirror.

* `mirror_locally` - (Required) Whether to copy the files locally or just
  references to the external files. Valid values are true or false.

* `name` - (Required) A name for the repo.

* `owners` - (Optional) List of Owners for authz_ownership.

* `proxy` - (Optional) Proxy to use for downloading the repo. This argument does
  not work on older versions of Cobbler.

* `rpm_list` - (Optional) List of specific RPMs to mirror.

## Attributes Reference

All of the above Optional attributes are also exported.
