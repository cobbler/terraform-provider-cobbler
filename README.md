# Cobbler Terraform Provider

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/d68c9aff2cd74b69afc9366ab4415f6a)](https://app.codacy.com/gh/cobbler/terraform-provider-cobbler/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/d68c9aff2cd74b69afc9366ab4415f6a)](https://app.codacy.com/gh/cobbler/terraform-provider-cobbler/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

The Cobbler provider is used to interact with a locally installed Cobbler service.\
The provider needs to be configured with the proper credentials before it can be used.

Original code by [Joe Topjian](https://github.com/jtopjian).

## Prerequisites

- [Terraform](https://terraform.io) 1.0 or above, **or** [OpenTofu](https://opentofu.org) 1.6 or above
- [Cobbler](https://cobbler.github.io/), release 4.0.0 (or higher)

## OpenTofu Support

This provider uses [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework) (protocol v6),
which is fully compatible with OpenTofu. You can use the provider with OpenTofu by referencing it from the
[OpenTofu Registry](https://search.opentofu.org/provider/cobbler/cobbler/latest).

```hcl
terraform {
  required_providers {
    cobbler = {
      source  = "cobbler/cobbler"
      version = "~> 6.0"
    }
  }
}
```

## Using the Provider

Full documentation can be found in the [`docs`](/docs) directory.

### Installation

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/cobbler/cobbler).

Make sure the file `variables.tf` contains the right version in the provider block:

```hcl
provider "cobbler" {
  version  = "~> 6.0.0"
  username = var.cobbler_username
  password = var.cobbler_password
  url      = var.cobbler_url
}
```

### Development

If you want to build from source, you can simply use `make` in the root of the repository.
#### Testing

To run the acceptance tests, type `make testacc`.  You will need [docker](https://docs.docker.com/get-docker/), 
[docker-compose](https://docs.docker.com/compose/install/) and xorriso installed.  Xorriso can be installed using:
`sudo apt-get install -y xorriso`, `sudo yum install xorriso -y`, or `sudo zypper install -y xorriso` depending on your
distro.  The Ubuntu 20.04 ISO will be downloaded idempotently to test importing a distro, this is < 1GB.  
