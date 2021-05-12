resource "cobbler_profile" "my_profile" {
  name        = "my_profile"
  distro      = "ubuntu-1804-x86_64"
  autoinstall = "default.ks"
}