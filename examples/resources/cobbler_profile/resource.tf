resource "cobbler_profile" "my_profile" {
  name        = "my_profile"
  distro      = "Ubuntu-2004-x86_64"
  autoinstall = "default.ks"
}