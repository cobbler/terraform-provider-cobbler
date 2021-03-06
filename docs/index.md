---
layout: ""
page_title: "Provider: Cobbler"
description: |-
  The Cobbler provider provides resources to interact with the Cobbler API.
---

# Cobbler Provider

The Cobbler provider is used to interact with a locally installed
[Cobbler](http://cobbler.github.io) service. The provider needs
to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Cobbler provider
provider "cobbler" {
  username = var.cobbler_username
  password = var.cobbler_password
  url      = var.cobbler_url
}

# Create a Cobbler Distro
resource "cobbler_distro" "ubuntu-1804-x86_64" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The username to the Cobbler service. This can
  also be specified with the `COBBLER_USERNAME` shell environment variable.

* `password` - (Required) The password to the Cobbler service. This can
  also be specified with the `COBBLER_PASSWORD` shell environment variable.

* `url` - (Required) The url to the Cobbler service. This can
  also be specified with the `COBBLER_URL` shell environment variable.

* `insecure` - (Optional) Ignore SSL certificate warnings and errors. This
  can also be specified with the `COBBLER_INSECURE` shell environment variable.

* `cacert_file` - (Optional) The path or contents of an SSL CA certificate.
  This can also be specified with the `COBBLER_CACERT_FILE` shell environment
  variable.
