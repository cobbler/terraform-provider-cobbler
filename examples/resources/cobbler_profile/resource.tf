resource "cobbler_distro" "ubuntu_2004" {
  name       = "Ubuntu-2004-x86_64"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "my_profile" {
  name        = "my_profile"
  distro      = cobbler_distro.ubuntu_2004.uid
  autoinstall = "default.ks"
}