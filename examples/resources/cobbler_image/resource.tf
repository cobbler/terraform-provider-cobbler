resource "cobbler_image" "Ubuntu-2004-x86_64" {
  name       = "foo"
  file       = "/var/www/cobbler/images/ubuntu-20.04-live-server-amd64.iso"
  breed      = "ubuntu"
  os_version = "focal"
  arch       = "x86_64"
  image_type = "iso"
}
