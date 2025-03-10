---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cobbler_profile Resource - terraform-provider-cobbler"
subcategory: ""
description: |-
  cobbler_profile manages a profile within Cobbler.
---

# cobbler_profile (Resource)

`cobbler_profile` manages a profile within Cobbler.

## Example Usage

```terraform
resource "cobbler_profile" "my_profile" {
  name        = "my_profile"
  distro      = "Ubuntu-2004-x86_64"
  autoinstall = "default.ks"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `distro` (String) Parent distribution.
- `name` (String) The name of the profile.

### Optional

- `autoinstall` (String) Template remote kickstarts or preseeds.
- `autoinstall_meta` (Map of String) Automatic installation template metadata, formerly Kickstart metadata.
- `autoinstall_meta_inherit` (Boolean) Signal that autoinstall_meta should be set to inherit from its parent
- `boot_files` (Map of String) Files copied into tftpboot beyond the kernel/initrd.
- `boot_files_inherit` (Boolean) Signal that boot_files should be set to inherit from its parent
- `comment` (String) Free form text description.
- `dhcp_tag` (String) DHCP tag.
- `enable_ipxe` (Boolean) Use iPXE instead of PXELINUX for advanced booting options.
- `enable_ipxe_inherit` (Boolean) Signal that enable_ipxe should be set to inherit from its parent
- `enable_menu` (Boolean) Enable a boot menu.
- `enable_menu_inherit` (Boolean) Signal that enable_menu should be set to inherit from its parent
- `fetchable_files` (Map of String) Templates for tftp or wget.
- `fetchable_files_inherit` (Boolean) Signal that fetchable_files should be set to inherit from its parent
- `kernel_options` (Map of String) Kernel options for the profile.
- `kernel_options_inherit` (Boolean) Signal that kernel_options should be set to inherit from its parent
- `kernel_options_post` (Map of String) Post install kernel options.
- `kernel_options_post_inherit` (Boolean) Signal that kernel_options_post should be set to inherit from its parent
- `mgmt_classes` (List of String) For external configuration management.
- `mgmt_classes_inherit` (Boolean) Signal that mgmt_classes should be set to inherit from its parent
- `mgmt_parameters` (Map of String) Parameters which will be handed to your management application (Must be a valid YAML dictionary).
- `mgmt_parameters_inherit` (Boolean) Signal that mgmt_parameters should be set to inherit from its parent
- `name_servers` (List of String) Name servers.
- `name_servers_inherit` (Boolean) Signal that name_servers should be set to inherit from its parent
- `name_servers_search` (List of String) Name server search settings.
- `name_servers_search_inherit` (Boolean) Signal that name_servers_search should be set to inherit from its parent
- `next_server_v4` (String) The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.
- `next_server_v6` (String) The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.
- `owners` (List of String) Owners list for authz_ownership.
- `owners_inherit` (Boolean) Signal that owners should be set to inherit from its parent
- `parent` (String) The parent this profile inherits settings from.
- `proxy` (String) Proxy URL.
- `repos` (List of String) Repos to auto-assign to this profile.
- `server` (String) The server-override for the profile.
- `template_files` (Map of String) File mappings for built-in config management.
- `template_files_inherit` (Boolean) Signal that template_files should be set to inherit from its parent
- `virt_auto_boot` (Boolean) Auto boot virtual machines.
- `virt_auto_boot_inherit` (Boolean) Signal that virt_auto_boot should be set to inherit from its parent
- `virt_bridge` (String) The bridge for virtual machines.
- `virt_cpus` (Number) The number of virtual CPUs
- `virt_disk_driver` (String) The virtual machine disk driver.
- `virt_file_size` (Number) The virtual machine file size.
- `virt_file_size_inherit` (Boolean) Signal that virt_file_size should be set to inherit from its parent
- `virt_path` (String) The virtual machine path.
- `virt_ram` (Number) The amount of RAM for the virtual machine.
- `virt_ram_inherit` (Boolean) Signal that virt_ram should be set to inherit from its parent
- `virt_type` (String) The type of virtual machine. Valid options are: xenpv, xenfv, qemu, kvm, vmware, openvz.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import cobbler_profile.foo foo
```
