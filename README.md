# Cobbler Terraform Provider

The Cobbler provider is used to interact with a locally installed Cobbler service.\
The provider needs to be configured with the proper credentials before it can be used.

Original code by [Joe Topjian](https://github.com/jtopjian).

## Prerequisites

- [Terraform](https://terraform.io), 0.14 and above
- [Cobbler](https://cobbler.github.io/), release 3.3.0 (or higher)

## Using the Provider

Full documentation can be found in the [`docs`](/docs) directory.

### Installation

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/cobbler/cobbler).

### Manual installation

You can download a pre-built binary from the [releases](https://github.com/cobbler/terraform-provider-cobbler/releases/)
 page.\
 These are built using [GoReleaser](https://goreleaser.com/) (the [configuration](.goreleaser.yml) is in the repo).

Download and add the pre-built binary for your system (Linux or macOS) to `~/.terraform.d/plugins/`.\
Replace `linux` with `darwin` for the macOS version.

```console
wget https://github.com/cobbler/terraform-provider-cobbler/releases/download/v2.0.3/terraform-provider-cobbler_2.0.3_linux_amd64.zip
unzip terraform-provider-cobbler_2.0.3_linux_amd64.zip
mkdir -p ~/.terraform.d/plugins/
mv terraform-provider-cobbler_v2.0.3 ~/.terraform.d/plugins/
chmod +x ~/.terraform.d/plugins/terraform-provider-cobbler_v2.0.3
```

Don't forget to run `terraform init` after installation of a new binary!

Make sure the file `variables.tf` contains the right version in the provider block:

```hcl
provider "cobbler" {
  version  = "~> 3.0.0"
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
