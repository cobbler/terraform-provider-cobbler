resource "cobbler_distro" "Ubuntu-2004-x86_64" {
  name       = "foo"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/linux"
  initrd     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/initrd.gz"
}