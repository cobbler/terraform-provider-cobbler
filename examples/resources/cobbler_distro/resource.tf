resource "cobbler_distro" "ubuntu-1804-x86_64" {
  name       = "foo"
  breed      = "ubuntu"
  os_version = "bionic"
  arch       = "x86_64"
  kernel     = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/linux"
  initrd     = "/var/www/cobbler/distro_mirror/Ubuntu-18.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
}