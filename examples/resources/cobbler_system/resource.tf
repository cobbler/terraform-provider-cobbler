resource "cobbler_distro" "ubuntu_2004" {
  name       = "Ubuntu-2004-x86_64"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  kernel     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/vmlinuz"
  initrd     = "/var/www/cobbler/distro_mirror/Ubuntu-20.04/install/initrd.gz"
}

resource "cobbler_profile" "my_profile" {
  name   = "my_profile"
  distro = cobbler_distro.ubuntu_2004.uid
}

resource "cobbler_system" "my_system" {
  name         = "my_system"
  profile      = cobbler_profile.my_profile.uid
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"
}

resource "cobbler_network_interface" "eth0" {
  name        = "eth0-${cobbler_system.my_system.name}"
  system      = cobbler_system.my_system.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "1.2.3.4"
    netmask = "255.255.255.0"
  }
}

resource "cobbler_network_interface" "eth1" {
  name        = "eth1-${cobbler_system.my_system.name}"
  system      = cobbler_system.my_system.uid
  mac_address = "aa:bb:cc:dd:ee:fa"
  static      = true
  ipv4 = {
    address = "1.2.3.5"
    netmask = "255.255.255.0"
  }
}