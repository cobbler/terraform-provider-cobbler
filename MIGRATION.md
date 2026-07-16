# Migration Guide: v5.0 → v6.0

Version 6.0 targets Cobbler 4.0.x via `cobblerclient` v1.0.0. The provider will
refuse to start against a Cobbler server older than 4.0.0 (the `Configure` step
fails with a `cobbler server too old` diagnostic). Users still on Cobbler 3.3.x
should stay on the v5.x line.

---

## Breaking Changes

### 1. `cobbler_snippet` and `cobbler_template_file` removed

Both resources are replaced by a single `cobbler_template` resource backed by the
new first-class Cobbler 4.0.0 `Template` item type. Body/name was replaced by
`content` + `name` + `template_type` + a `uri = { schema, path }` block.

**Before (v5):**

```hcl
resource "cobbler_template_file" "preseed" {
  name = "ubuntu.preseed"
  body = file("ubuntu.preseed")
}
```

**After (v6):**

```hcl
resource "cobbler_template" "preseed" {
  name          = "ubuntu.preseed"
  template_type = "jinja"
  uri = {
    schema = "file"
    path   = "ubuntu.preseed"
  }
  content = file("ubuntu.preseed")
}
```

Re-import each snippet/template_file as a `cobbler_template`:

```bash
terraform state rm cobbler_template_file.preseed
terraform import cobbler_template.preseed ubuntu.preseed
```

### 2. `cobbler_system.interface` map removed; use `cobbler_network_interface`

Network interfaces are now first-class Cobbler 4.0.0 items, not nested attributes
of `cobbler_system`. Each entry of the old `interface` map becomes its own
`cobbler_network_interface` resource. Network interfaces are a flat, top-level
Cobbler collection, so `name` must now be globally unique across all systems
(Cobbler's `validate_obj_name` rejects `@`, so the old Cobbler 3.x
`<ifname>@<systemname>` convention doesn't work here - pick any unique name,
e.g. `<ifname>-<systemname>`), and the per-interface IPv4/IPv6/DNS settings are
nested objects.

**Before (v5):**

```hcl
resource "cobbler_system" "foo" {
  name    = "foo"
  profile = "ubuntu"
  interface = {
    "eth0" = {
      mac_address = "aa:bb:cc:dd:ee:ff"
      static      = true
      ip_address  = "10.0.0.5"
      netmask     = "255.255.255.0"
    }
  }
}
```

**After (v6):**

```hcl
resource "cobbler_system" "foo" {
  name    = "foo"
  profile = cobbler_profile.ubuntu.uid
}

resource "cobbler_network_interface" "foo_eth0" {
  name        = "eth0-${cobbler_system.foo.name}"
  system      = cobbler_system.foo.uid
  mac_address = "aa:bb:cc:dd:ee:ff"
  static      = true
  ipv4 = {
    address = "10.0.0.5"
    netmask = "255.255.255.0"
  }
}
```

To migrate state:

```bash
terraform state rm cobbler_system.foo
terraform import cobbler_system.foo foo
terraform import cobbler_network_interface.foo_eth0 eth0-foo
```

### 3. `mgmt_classes` and `mgmt_parameters` attributes removed

The `MgmtClass` item type was removed from Cobbler 4.0.0 server-side, taking the
`mgmt_classes` and `mgmt_parameters` fields with it. Remove these attributes from
your HCL for `cobbler_distro`, `cobbler_image`, `cobbler_menu`, `cobbler_profile`,
and `cobbler_system`.

### 4. New `cobbler_system.uid` computed attribute

A `uid` computed attribute has been added so that `cobbler_network_interface.system`
can reference the parent system by its server-assigned UID.

### 5. New resources

- `cobbler_template` — replaces snippet + template_file
- `cobbler_network_interface` — replaces the `interface` map on `cobbler_system`
- `cobbler_distro_group`, `cobbler_profile_group`, `cobbler_system_group` —
  named collections of distros/profiles/systems for bulk operations

Each new resource has a matching data source.

### 6. Server requirement

Cobbler ≥ 4.0.0 required. The provider's `Configure` step will fail with
`cobbler server too old` if the connected server reports a lower version.

### 7. Cross-reference fields now take a UID, not a name

`cobbler_profile.distro`, `cobbler_profile.parent`, `cobbler_system.profile`,
`cobbler_system.image`, and `cobbler_image.menu` now take the referenced
resource's Cobbler UID instead of its name. Cobbler's own API has always been
UID-only for these fields server-side; the provider previously relied on
`cobblerclient` translating names to UIDs and back, which meant renaming the
referenced object produced a permanent, un-satisfiable plan diff (and could
force an unwanted resource recreation) since the name in your config no
longer matched the UID Cobbler actually stored. Referencing the UID directly
avoids that entirely.

New computed `uid` attributes were added to `cobbler_distro`, `cobbler_profile`,
`cobbler_image`, and `cobbler_menu` for this (`cobbler_system` already had one,
added for `cobbler_network_interface.system`).

**Before (v5):**

```hcl
resource "cobbler_distro" "ubuntu" {
  name = "ubuntu"
  # ...
}

resource "cobbler_profile" "foo" {
  name   = "foo"
  distro = cobbler_distro.ubuntu.name
}
```

**After (v6):**

```hcl
resource "cobbler_distro" "ubuntu" {
  name = "ubuntu"
  # ...
}

resource "cobbler_profile" "foo" {
  name   = "foo"
  distro = cobbler_distro.ubuntu.uid
}
```

**State impact:** existing state stores the referenced object's name in these
fields. On the first `plan`/`apply` after upgrading, Terraform will show a
diff changing these fields from the name to the UID — this is expected and
one-time; apply it to bring state in line with the new UID-based format.

---

# Migration Guide: v4.x → v5.0

Version 5.0 migrates the provider from `terraform-plugin-sdk/v2` to `terraform-plugin-framework`. This is a
**protocol v6** provider and requires Terraform ≥ 1.0 or OpenTofu ≥ 1.6.

---

## Breaking Changes

### 1. Inheritable fields: `foo` + `foo_inherit` → nested `foo` object

Every field that Cobbler can inherit from a parent profile/distro has changed from a flat pair of attributes to a single
nested object:

**Before:**

```hcl
resource "cobbler_profile" "example" {
  name_servers       = ["8.8.8.8"]
  name_servers_inherit = false
}
```

**After:**

```hcl
resource "cobbler_profile" "example" {
  name_servers = {
    value     = ["8.8.8.8"]
    inherited = false
  }
}
```

To use the Cobbler-inherited value instead of overriding it:

```hcl
name_servers = {
  value     = []
  inherited = true
}
```

This change applies to all inheritable fields across `cobbler_distro`, `cobbler_repo`, `cobbler_profile`, and
`cobbler_system`.

### 2. `cobbler_system` interface: TypeSet blocks → map attribute

The `interface` block has changed from a set of blocks (keyed by `name`) to a map attribute keyed by the interface name.

**Before:**

```hcl
resource "cobbler_system" "example" {
  interface {
    name        = "eth0"
    mac_address = "aa:bb:cc:dd:ee:ff"
    ip_address  = "1.2.3.4"
    netmask     = "255.255.255.0"
    static      = true
  }
}
```

**After:**

```hcl
resource "cobbler_system" "example" {
  interface = {
    "eth0" = {
      mac_address = "aa:bb:cc:dd:ee:ff"
      ip_address  = "1.2.3.4"
      netmask     = "255.255.255.0"
      static      = true
    }
  }
}
```

Note: the `name` sub-attribute is gone — the map key IS the interface name.

### 3. `cobbler_repo.environment`: string → map(string)

The `environment` attribute on `cobbler_repo` changed from a flat string to a map of key-value pairs.

**Before (v3.x):**

```hcl
resource "cobbler_repo" "example" {
  name    = "my-repo"
  mirror  = "http://example.com/repo"
  breed   = "yum"
  environment = "VAR1=value1 VAR2=value2"
}
```

**After (v4.0):**

```hcl
resource "cobbler_repo" "example" {
  name    = "my-repo"
  mirror  = "http://example.com/repo"
  breed   = "yum"
  environment = {
    VAR1 = "value1"
    VAR2 = "value2"
  }
}
```

**State impact:** If you have existing `cobbler_repo` resources in Terraform state that have `environment` set, the
state will contain the old string encoding. Terraform will fail to deserialize the state after upgrading to v5.0.
You must remove those resources from state and re-import them (see [State Migration Steps](#state-migration-steps)
below).

### 4. Read-only data sources added

Data sources (`data "cobbler_*"`) now exist for all six resource types. These are new and have no migration impact
unless you have existing `data` blocks that now conflict.

---

## State Migration Steps

Terraform state is provider-SDK-encoded. Because the schema has changed significantly, the safest upgrade path is to
remove existing resources from state and re-import them.

```bash
# 1. List all cobbler resources in state
terraform state list | grep '^cobbler_'

# 2. Remove each resource from state (repeat for every resource)
terraform state rm cobbler_system.example
terraform state rm cobbler_profile.example
# ... etc.

# 3. Update your HCL to the new syntax (see above)

# 4. Re-import each resource by its name (Cobbler object name)
terraform import cobbler_system.example my_system_name
terraform import cobbler_profile.example my_profile_name
# ... etc.

# 5. Run plan to confirm no diff
terraform plan
```

### `cobbler_repo` resources with `environment` set

If you have `cobbler_repo` resources that used the `environment` attribute under v4.x, those resources **must** be
removed from state and re-imported after upgrading. The attribute type changed from a string to `map(string)`, and no
automatic state upgrade is performed.

```bash
# Identify affected repo resources
terraform state list | grep '^cobbler_repo\.'

# For each repo that had environment set, remove it from state and re-import
terraform state rm cobbler_repo.my_repo
# Update HCL to use the new map syntax (see Breaking Change #3 above)
terraform import cobbler_repo.my_repo my_repo_name

# Confirm no diff
terraform plan
```

---

## OpenTofu Usage

This provider (v5.0+) uses protocol v6 and works with OpenTofu out of the box.

```hcl
terraform {
  required_providers {
    cobbler = {
      source  = "cobbler/cobbler"
      version = "~> 5.0"
    }
  }
}
```

Run `tofu init` instead of `terraform init`.
