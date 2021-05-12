resource "cobbler_system" "my_system" {
  name         = "my_system"
  profile      = "my_profile"
  name_servers = ["8.8.8.8", "8.8.4.4"]
  comment      = "I'm a system"

  interface {
    name        = "eth0"
    mac_address = "aa:bb:cc:dd:ee:ff"
    static      = true
    ip_address  = "1.2.3.4"
    netmask     = "255.255.255.0"
  }

  interface {
    name        = "eth1"
    mac_address = "aa:bb:cc:dd:ee:fa"
    static      = true
    ip_address  = "1.2.3.5"
    netmask     = "255.255.255.0"
  }
}