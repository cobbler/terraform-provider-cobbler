resource "cobbler_repo" "my_repo" {
  name           = "my_repo"
  breed          = "apt"
  arch           = "x86_64"
  apt_components = ["main"]
  apt_dists      = ["focal"]
  mirror         = "http://us.archive.ubuntu.com/ubuntu/"
}