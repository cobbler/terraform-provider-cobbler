#!/bin/bash

set -e

# This script assumes Ubuntu 18.04 is being used.
# It will create a standard Cobbler environment that can be used for acceptance testing.

# With this environment spun up, the config should be:
#  COBBLER_URL=http://127.0.0.1:25151
#  COBBLER_USERNAME=cobbler
#  COBBLER_PASSWORD=cobbler

sudo apt-get update
sudo apt-get install -y build-essential git mercurial

cd
echo 'export PATH=$PATH:$HOME/terraform:$HOME/go/bin' >> ~/.bashrc
export PATH=$PATH:$HOME/terraform:$HOME/go/bin

sudo wget -O /usr/local/bin/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
sudo chmod +x /usr/local/bin/gimme
/usr/local/bin/gimme 1.12 >> ~/.bashrc
eval "$(/usr/local/bin/gimme 1.12)"

mkdir ~/go
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc
export GOPATH=$HOME/go
source ~/.bashrc

#git clone https://github.com/terraform-providers/terraform-provider-cobbler
git clone https://github.com/wearespindle/terraform-provider-cobbler

# Cobbler
sudo apt-get install -y cobbler cobbler-web debmirror dnsmasq

sudo tee /etc/cobbler/modules.conf <<EOF
[authentication]
module = authentication.configfile
[authorization]
module = authorization.allowall
[dns]
module = managers.dnsmasq
[dhcp]
module = managers.dnsmasq
[tftpd]
module = managers.in_tftpd
EOF

sudo tee /etc/cobbler/dnsmasq.template <<EOF
dhcp-range = 192.168.255.200,192.168.255.250
server = 8.8.8.8
read-ethers
addn-hosts = /var/lib/cobbler/cobbler_hosts

dhcp-option=3,\$next_server
dhcp-lease-max=1000
dhcp-authoritative
dhcp-boot=pxelinux.0
dhcp-boot=net:normalarch,pxelinux.0
dhcp-boot=net:ia64,\$elilo

\$insert_cobbler_system_definitions
EOF

sudo sed -i -e 's/^manage_dhcp: 0/manage_dhcp: 1/' /etc/cobbler/settings
sudo sed -i -e 's/^manage_dns: 0/manage_dns: 1/' /etc/cobbler/settings
sudo sed -i -e 's/^next_server:.*/next_server: 127.0.0.1/' /etc/cobbler/settings
sudo sed -i -e 's/^server:.*/server: 127.0.0.1/' /etc/cobbler/settings

# User: cobbler / Pass: cobbler
sudo tee /etc/cobbler/users.digest <<EOF
cobbler:Cobbler:a2d6bae81669d707b72c0bd9806e01f3
EOF

sudo systemctl daemon-reload
sudo systemctl restart httpd
sudo systemctl stop cobblerd
sleep 2
sudo systemctl start cobblerd
sleep 5
sudo cobbler get-loaders
sudo cobbler sync

# Import a recent Ubuntu distro
cd /tmp
wget http://cdimage.ubuntu.com/ubuntu/releases/18.04/release/ubuntu-18.04.4-server-amd64.iso
sudo mount -o loop ubuntu-18.04.4-server-amd64.iso /mnt
sudo cobbler import --name Ubuntu-18.04 --breed ubuntu --path /mnt

# Create a file with the cobbler credential environment variables
cat > ~/cobblerc <<EOF
export COBBLER_USERNAME="cobbler"
export COBBLER_PASSWORD="cobbler"
export COBBLER_URL="http://localhost:25151"
EOF
