package cobbler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"sync"

	cobbler "github.com/cobbler/cobblerclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemSyncLock sync.Mutex

func resourceSystem() *schema.Resource {
	return &schema.Resource{
		Description:   "`cobbler_system` manages a system within Cobbler.",
		CreateContext: resourceSystemCreate,
		ReadContext:   resourceSystemRead,
		UpdateContext: resourceSystemUpdate,
		DeleteContext: resourceSystemDelete,

		Schema: map[string]*schema.Schema{
			"autoinstall": {
				Description: "Template remote kickstarts or preseeds.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"autoinstall_meta": {
				Description: "Automatic installation template metadata, formerly Kickstart metadata.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"boot_files": {
				Description: "Files copied into tftpboot beyond the kernel/initrd.",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"boot_loaders": {
				Description: "Must be either `grub`, `pxe`, or `ipxe`.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},

			"comment": {
				Description: "Free form text description.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"enable_gpxe": {
				Description: "Use gPXE instead of PXELINUX for advanced booting options.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},

			"fetchable_files": {
				Description: "Templates for tftp or wget.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},

			"gateway": {
				Description: "Network gateway.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"hostname": {
				Description: "Hostname of the system.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"image": {
				Description: "Parent image (if no profile is used).",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"interface": {
				Description: "The `interface` Block Set.",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The device name of the interface. ex: `eth0`.",
							Type:        schema.TypeString,
							Required:    true,
						},

						"cnames": {
							Description: "Canonical name records.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"dhcp_tag": {
							Description: "DHCP tag.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"dns_name": {
							Description: "DNS name.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"bonding_opts": {
							Description: "Options for bonded interfaces.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"bridge_opts": {
							Description: "Options for bridge interfaces.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"gateway": {
							Description: "Per-interface gateway.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"interface_type": {
							// TODO: Update list of interface types
							Description:  "The type of interface: NA, master, slave, bond, bond_slave, bridge, bridge_slave, bonded_bridge_slave, infiniband, bmc",
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NA",
							ValidateFunc: validation.StringInSlice([]string{"NA", "master", "slave", "bond", "bond_slave", "bridge", "bridge_slave", "bonded_bridge_slave", "infiniband", "bmc"}, false),
						},
						"interface_master": {
							Description: "The master interface when slave.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ip_address": {
							Description: "The IP address of the interface.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ipv6_address": {
							Description: "The IPv6 address of the interface.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ipv6_secondaries": {
							Description: "IPv6 secondaries.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv6_mtu": {
							Description: "The MTU of the IPv6 address.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"ipv6_static_routes": {
							Description: "Static routes for the IPv6 interface.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv6_default_gateway": {
							Description: "The default gateawy for the IPv6 address / interface.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"mac_address": {
							Description: "The MAC address of the interface.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"management": {
							Description: "Whether this interface is a management interface.",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"netmask": {
							Description: "The IPv4 netmask of the interface.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"static": {
							Description: "Whether the interface should be static or DHCP.",
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"static_routes": {
							Description: "Static routes for the interface.",
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"virt_bridge": {
							Description: "The virtual bridge to attach to.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
					},
				},
				Set: resourceSystemInterfaceHash,
			},
			"ipv6_default_device": {
				Description: "IPv6 default device.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"kernel_options": {
				Description: "Kernel options. ex: `selinux=permissive`.",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"kernel_options_post": {
				Description: "Kernel options (post install).",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"mgmt_classes": {
				Description: "For external configuration management.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mgmt_parameters": {
				Description: "Parameters which will be handed to your management application (Must be a valid YAML dictionary).",
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "The name of the system.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name_servers_search": {
				Description: "Name server search settings.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name_servers": {
				Description: "Name servers.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"netboot_enabled": {
				Description: "(Re)install this machine at next boot.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"next_server_v4": {
				Description: "The next_server_v4 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"next_server_v6": {
				Description: "The next_server_v6 option is used for DHCP/PXE as the IP of the TFTP server from which network boot files are downloaded. Usually, this will be the same IP as the server setting.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"owners": {
				Description: "Owners list for authz_ownership.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"power_address": {
				Description: "Power management address.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"power_id": {
				Description: "Usually a plug number or blade name if power type requires it.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"power_pass": {
				Description: "Power management password.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"power_type": {
				Description: "Power management type.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"power_user": {
				Description: "Power management user.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"profile": {
				Description: "Parent profile.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"proxy": {
				Description: "Proxy URL.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"status": {
				Description: "System status (development, testing, acceptance, production).",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"template_files": {
				Description: "File mappings for built-in config management.",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
			},
			"virt_auto_boot": {
				Description: "Auto boot virtual machines.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"virt_file_size": {
				Description: "The virtual machine file size.",
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
			},
			"virt_cpus": {
				Description: "The number of virtual CPUs",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"virt_type": {
				Description: "The type of virtual machine. Valid options are: xenpv, xenfv, qemu, kvm, vmware, openvz.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"virt_path": {
				Description: "The virtual machine path.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"virt_pxe_boot": {
				Description: "Use PXE to build this virtual machine.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"virt_ram": {
				Description: "The amount of RAM for the virtual machine.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"virt_disk_driver": {
				Description: "The virtual machine disk driver.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceSystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	config := meta.(*Config)

	// Create a cobblerclient.System struct
	system, err := buildSystem(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to create the System
	tflog.Debug(ctx, "Cobbler System: Create Options", map[string]interface{}{
		"options": structs.Map(system),
	})
	newSystem, err := config.cobblerClient.CreateSystem(system)
	if err != nil {
		return diag.Errorf("Cobbler System: Error Creating: %s", err)
	}

	// Build cobblerclient.Interface structs
	interfaces := buildSystemInterfaces(d.Get("interface").(*schema.Set))

	// Add each interface to the system
	for interfaceName, interfaceInfo := range interfaces {
		tflog.Debug(ctx, "Cobbler System Interface", map[string]interface{}{
			"interface": interfaceName,
			"options":   structs.Map(interfaceInfo),
		})
		if err = newSystem.CreateInterface(interfaceName, interfaceInfo); err != nil {
			return diag.Errorf("Cobbler System: Error adding Interface %s to %s: %s", interfaceName, newSystem.Name, err)
		}
	}

	tflog.Debug(ctx, "Cobbler System: Created System", map[string]interface{}{
		"system": structs.Map(newSystem),
	})
	d.SetId(newSystem.Name)

	tflog.Debug(ctx, "Cobbler System: syncing system")
	if err = config.cobblerClient.Sync(); err != nil {
		return diag.Errorf("Cobbler System: Error syncing system: %s", err)
	}

	return resourceSystemRead(ctx, d, meta)
}

func resourceSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	tflog.Debug(ctx, "Reading Cobbler system", map[string]interface{}{
		"system": d.Id(),
	})

	// Retrieve the system entry from Cobbler
	system, err := config.cobblerClient.GetSystem(d.Id(), false, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Warn(ctx, "Cobbler System not found, removing from state", map[string]interface{}{
				"system": d.Id(),
			})
			d.SetId("")
			return nil
		}

		return diag.Errorf("Cobbler System: Error Reading (%s): %s", d.Id(), err)
	}

	// Set all fields
	err = d.Set("boot_files", system.BootFiles.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("boot_loaders", system.BootLoaders.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("comment", system.Comment)
	if err != nil {
		return diag.FromErr(err)
	}
	// TODO: enable_ipxe
	err = d.Set("enable_gpxe", system.EnableIPXE.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("fetchable_files", system.FetchableFiles.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("gateway", system.Gateway)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("hostname", system.Hostname)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("image", system.Image)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("ipv6_default_device", system.IPv6DefaultDevice)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options", system.KernelOptions.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("kernel_options_post", system.KernelOptionsPost.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("autoinstall_meta", system.AutoinstallMeta.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("mgmt_classes", system.MgmtClasses.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	if system.MgmtParameters.IsInherited {
		err = d.Set("mgmt_parameters", make(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = d.Set("mgmt_parameters", system.MgmtParameters.Data)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = d.Set("name", system.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name_servers_search", system.NameServersSearch)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("name_servers", system.NameServers)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("netboot_enabled", system.NetbootEnabled)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("next_server_v4", system.NextServerv4)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("next_server_v6", system.NextServerv6)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("owners", system.Owners.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("power_address", system.PowerAddress)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("power_id", system.PowerID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("power_pass", system.PowerPass)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("power_type", system.PowerType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("power_user", system.PowerUser)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("profile", system.Profile)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("proxy", system.Proxy)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("status", system.Status)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("template_files", system.TemplateFiles.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_auto_boot", system.VirtAutoBoot.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	if !system.VirtFileSize.IsInherited {
		err = d.Set("virt_file_size", system.VirtFileSize.Data)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = d.Set("virt_cpus", system.VirtCPUs.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_type", system.VirtType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_path", system.VirtPath)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_pxe_boot", system.VirtPXEBoot)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_ram", system.VirtRAM.Data)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("virt_disk_driver", system.VirtDiskDriver)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get all interfaces that the System has
	allInterfaces, err := system.GetInterfaces()
	if err != nil {
		return diag.Errorf("Cobbler System %s: Error getting interfaces: %s", system.Name, err)
	}

	// Build a generic map array with the interface attributes
	var systemInterfaces []map[string]interface{}
	for interfaceName, interfaceInfo := range allInterfaces {
		tflog.Debug(ctx, "Cobbler System Interface", map[string]interface{}{
			"interface": interfaceName,
			"options":   structs.Map(interfaceInfo),
		})
		iface := make(map[string]interface{})
		iface["name"] = interfaceName
		iface["cnames"] = interfaceInfo.CNAMEs
		iface["dhcp_tag"] = interfaceInfo.DHCPTag
		iface["dns_name"] = interfaceInfo.DNSName
		iface["bonding_opts"] = interfaceInfo.BondingOpts
		iface["bridge_opts"] = interfaceInfo.BridgeOpts
		iface["gateway"] = interfaceInfo.Gateway
		iface["interface_type"] = interfaceInfo.InterfaceType
		iface["interface_master"] = interfaceInfo.InterfaceMaster
		iface["ip_address"] = interfaceInfo.IPAddress
		iface["ipv6_address"] = interfaceInfo.IPv6Address
		iface["ipv6_secondaries"] = interfaceInfo.IPv6Secondaries
		iface["ipv6_mtu"] = interfaceInfo.IPv6MTU
		iface["ipv6_static_routes"] = interfaceInfo.IPv6StaticRoutes
		iface["ipv6_default_gateway"] = interfaceInfo.IPv6DefaultGateway
		iface["mac_address"] = interfaceInfo.MACAddress
		iface["management"] = interfaceInfo.Management
		iface["netmask"] = interfaceInfo.Netmask
		iface["static"] = interfaceInfo.Static
		iface["static_routes"] = interfaceInfo.StaticRoutes
		iface["virt_bridge"] = interfaceInfo.VirtBridge
		systemInterfaces = append(systemInterfaces, iface)
	}

	err = d.Set("interface", systemInterfaces)
	if err != nil {
		return diag.Errorf("Cobbler System %s: Error appending interface to : %s", system.Name, err)
	}

	return nil
}

func resourceSystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	systemSyncLock.Lock()
	defer systemSyncLock.Unlock()

	config := meta.(*Config)

	// Retrieve the existing system entry from Cobbler
	system, err := config.cobblerClient.GetSystem(d.Id(), false, false)
	if err != nil {
		return diag.Errorf("Cobbler System: Error Reading (%s): %s", d.Id(), err)
	}

	// Get a list of the old interfaces
	currentInterfaces, err := system.GetInterfaces()
	if err != nil {
		return diag.Errorf("error getting interfaces: %s", err)
	}
	interfaceMap := make(map[string]map[string]interface{})
	for interfaceName, interfaceInfo := range currentInterfaces {
		interfaceMap[interfaceName] = structs.Map(interfaceInfo)
	}
	tflog.Debug(ctx, "Cobbler System Interfaces", map[string]interface{}{
		"interfaces": interfaceMap,
	})

	// Create a new cobblerclient.System struct with the new information
	newSystem, err := buildSystem(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to update the system with new information
	tflog.Debug(ctx, "Cobbler System: Updating System with options", map[string]interface{}{
		"system":  d.Id(),
		"options": structs.Map(system),
	})
	err = config.cobblerClient.UpdateSystem(&newSystem)
	if err != nil {
		return diag.Errorf("Cobbler System: Error Updating (%s): %s", d.Id(), err)
	}

	if d.HasChange("interface") {
		oldInterfaces, newInterfaces := d.GetChange("interface")
		oldInterfacesSet := oldInterfaces.(*schema.Set)
		newInterfacesSet := newInterfaces.(*schema.Set)
		interfacesToRemove := oldInterfacesSet.Difference(newInterfacesSet)

		oldIfaces := buildSystemInterfaces(interfacesToRemove)
		newIfaces := buildSystemInterfaces(newInterfacesSet)

		for interfaceName, interfaceInfo := range oldIfaces {
			if _, ok := newIfaces[interfaceName]; !ok {
				// Interface does not exist in the new set, so it has been removed from terraform.
				tflog.Debug(ctx, "Cobbler System: Deleting Interface", map[string]interface{}{
					"interface": interfaceName,
					"options":   structs.Map(interfaceInfo),
				})

				if err := system.DeleteInterface(interfaceName); err != nil {
					return diag.Errorf("Cobbler System: Error deleting Interface %s to %s: %s", interfaceName, system.Name, err)
				}
			}
		}

		// Modify interfaces that have changed
		for interfaceName, interfaceInfo := range newIfaces {
			tflog.Debug(ctx, "Cobbler System: New Interface", map[string]interface{}{
				"interface": interfaceName,
				"options":   structs.Map(interfaceInfo),
			})

			if err := system.CreateInterface(interfaceName, interfaceInfo); err != nil {
				return diag.Errorf("Cobbler System: Error adding Interface %s to %s: %s", interfaceName, system.Name, err)
			}
		}
	}

	tflog.Debug(ctx, "Cobbler System: syncing system")
	if err := config.cobblerClient.Sync(); err != nil {
		return diag.Errorf("Cobbler System: Error syncing system: %s", err)
	}

	return resourceSystemRead(ctx, d, meta)
}

func resourceSystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Attempt to delete the system
	if err := config.cobblerClient.DeleteSystem(d.Id()); err != nil {
		return diag.Errorf("Cobbler System: Error Deleting (%s): %s", d.Id(), err)
	}

	return nil
}

// buildSystem builds a cobblerclient.System out of the Terraform attributes.
func buildSystem(d *schema.ResourceData) (cobbler.System, error) {
	mgmtClasses, err := GetStringSlice(d, "mgmt_classes")
	if err != nil {
		return cobbler.System{}, err
	}
	mgmtParameters, err := GetInterfaceMap(d, "mgmt_parameters")
	if err != nil {
		return cobbler.System{}, err
	}
	nameServersSearch, err := GetStringSlice(d, "name_servers_search")
	if err != nil {
		return cobbler.System{}, err
	}
	nameServers, err := GetStringSlice(d, "name_servers")
	if err != nil {
		return cobbler.System{}, err
	}
	owners, err := GetStringSlice(d, "owners")
	if err != nil {
		return cobbler.System{}, err
	}
	bootFiles, err := GetInterfaceMap(d, "boot_files")
	if err != nil {
		return cobbler.System{}, err
	}
	fetchableFiles, err := GetInterfaceMap(d, "fetchable_files")
	if err != nil {
		return cobbler.System{}, err
	}
	kernelOptions, err := GetInterfaceMap(d, "kernel_options")
	if err != nil {
		return cobbler.System{}, err
	}
	kernelOptionsPost, err := GetInterfaceMap(d, "kernel_options_post")
	if err != nil {
		return cobbler.System{}, err
	}
	templateFiles, err := GetInterfaceMap(d, "template_files")
	if err != nil {
		return cobbler.System{}, err
	}
	autoinstallMeta, err := GetInterfaceMap(d, "autoinstall_meta")
	if err != nil {
		return cobbler.System{}, err
	}
	bootLoaders, err := GetStringSlice(d, "boot_loaders")
	if err != nil {
		return cobbler.System{}, err
	}

	system := cobbler.NewSystem()
	system.Autoinstall = d.Get("autoinstall").(string)
	system.AutoinstallMeta = cobbler.Value[map[string]interface{}]{
		Data:        autoinstallMeta,
		IsInherited: false,
	}
	system.BootFiles = cobbler.Value[map[string]interface{}]{
		Data:        bootFiles,
		IsInherited: false,
	}
	system.BootLoaders = cobbler.Value[[]string]{
		Data:        bootLoaders,
		IsInherited: false,
	}
	system.Comment = d.Get("comment").(string)
	// TODO: enable_ipxe
	system.EnableIPXE = cobbler.Value[bool]{
		Data:        d.Get("enable_gpxe").(bool),
		IsInherited: false,
	}
	system.FetchableFiles = cobbler.Value[map[string]interface{}]{
		Data:        fetchableFiles,
		IsInherited: false,
	}
	system.Gateway = d.Get("gateway").(string)
	system.Hostname = d.Get("hostname").(string)
	system.Image = d.Get("image").(string)
	system.IPv6DefaultDevice = d.Get("ipv6_default_device").(string)
	system.KernelOptions = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptions,
		IsInherited: false,
	}
	system.KernelOptionsPost = cobbler.Value[map[string]interface{}]{
		Data:        kernelOptionsPost,
		IsInherited: false,
	}
	system.MgmtClasses = cobbler.Value[[]string]{
		Data:        mgmtClasses,
		IsInherited: false,
	}
	system.MgmtParameters = cobbler.Value[map[string]interface{}]{
		Data:        mgmtParameters,
		IsInherited: false,
	}
	system.Name = d.Get("name").(string)
	system.NameServersSearch = nameServersSearch
	system.NameServers = nameServers
	system.NetbootEnabled = d.Get("netboot_enabled").(bool)
	system.NextServerv4 = d.Get("next_server_v4").(string)
	system.NextServerv6 = d.Get("next_server_v6").(string)
	system.Owners = cobbler.Value[[]string]{
		Data:        owners,
		IsInherited: false,
	}
	system.PowerAddress = d.Get("power_address").(string)
	system.PowerID = d.Get("power_id").(string)
	system.PowerPass = d.Get("power_pass").(string)
	system.PowerType = d.Get("power_type").(string)
	system.PowerUser = d.Get("power_user").(string)
	system.Profile = d.Get("profile").(string)
	system.Proxy = d.Get("proxy").(string)
	system.Status = d.Get("status").(string)
	system.TemplateFiles = cobbler.Value[map[string]interface{}]{
		Data:        templateFiles,
		IsInherited: false,
	}
	system.VirtAutoBoot = cobbler.Value[bool]{
		Data:        d.Get("virt_auto_boot").(bool),
		IsInherited: false,
	}
	system.VirtFileSize = cobbler.Value[float64]{
		Data:        d.Get("virt_file_size").(float64),
		IsInherited: false,
	}
	system.VirtCPUs = cobbler.Value[int]{
		Data:        d.Get("virt_cpus").(int),
		IsInherited: false,
	}
	system.VirtType = d.Get("virt_type").(string)
	system.VirtPath = d.Get("virt_path").(string)
	system.VirtPXEBoot = d.Get("virt_pxe_boot").(bool)
	system.VirtRAM = cobbler.Value[int]{
		Data:        d.Get("virt_ram").(int),
		IsInherited: false,
	}
	system.VirtDiskDriver = d.Get("virt_disk_driver").(string)

	return system, err
}

// buildSystemInterface builds a cobblerclient.Interface out of the Terraform attributes.
func buildSystemInterfaces(systemInterfaces *schema.Set) cobbler.Interfaces {
	interfaces := make(cobbler.Interfaces)
	rawInterfaces := systemInterfaces.List()
	for _, rawInterface := range rawInterfaces {
		rawInterfaceMap := rawInterface.(map[string]interface{})

		cnames := []string{}
		for _, i := range rawInterfaceMap["cnames"].([]interface{}) {
			cnames = append(cnames, i.(string))
		}

		ipv6Secondaries := []string{}
		for _, i := range rawInterfaceMap["ipv6_secondaries"].([]interface{}) {
			ipv6Secondaries = append(ipv6Secondaries, i.(string))
		}

		ipv6StaticRoutes := []string{}
		for _, i := range rawInterfaceMap["ipv6_static_routes"].([]interface{}) {
			ipv6StaticRoutes = append(ipv6StaticRoutes, i.(string))
		}

		staticRoutes := []string{}
		for _, i := range rawInterfaceMap["static_routes"].([]interface{}) {
			staticRoutes = append(staticRoutes, i.(string))
		}

		interfaceName := rawInterfaceMap["name"].(string)
		interfaces[interfaceName] = cobbler.Interface{
			CNAMEs:             cnames,
			DHCPTag:            rawInterfaceMap["dhcp_tag"].(string),
			DNSName:            rawInterfaceMap["dns_name"].(string),
			BondingOpts:        rawInterfaceMap["bonding_opts"].(string),
			BridgeOpts:         rawInterfaceMap["bridge_opts"].(string),
			Gateway:            rawInterfaceMap["gateway"].(string),
			InterfaceType:      rawInterfaceMap["interface_type"].(string),
			InterfaceMaster:    rawInterfaceMap["interface_master"].(string),
			IPAddress:          rawInterfaceMap["ip_address"].(string),
			IPv6Address:        rawInterfaceMap["ipv6_address"].(string),
			IPv6Secondaries:    ipv6Secondaries,
			IPv6MTU:            rawInterfaceMap["ipv6_mtu"].(string),
			IPv6StaticRoutes:   ipv6StaticRoutes,
			IPv6DefaultGateway: rawInterfaceMap["ipv6_default_gateway"].(string),
			MACAddress:         rawInterfaceMap["mac_address"].(string),
			Management:         rawInterfaceMap["management"].(bool),
			Netmask:            rawInterfaceMap["netmask"].(string),
			Static:             rawInterfaceMap["static"].(bool),
			StaticRoutes:       staticRoutes,
			VirtBridge:         rawInterfaceMap["virt_bridge"].(string),
		}
	}

	return interfaces
}

func resourceSystemInterfaceHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(m["name"].(string))

	if v, ok := m["mac_address"]; ok {
		buf.WriteString(fmt.Sprintf("%v-", v.(string)))
	}

	hash := String(buf.String())
	log.Printf("[DEBUG] Interface %s: Calculated hash %v", m["name"].(string), hash)
	return hash
}
