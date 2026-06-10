# Migration Guide: v3.x → v4.0

Version 4.0 migrates the provider from `terraform-plugin-sdk/v2` to `terraform-plugin-framework`. This is a
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
state will contain the old string encoding. Terraform will fail to deserialize the state after upgrading to v4.0.
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

If you have `cobbler_repo` resources that used the `environment` attribute under v3.x, those resources **must** be
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

This provider (v4.0+) uses protocol v6 and works with OpenTofu out of the box.

```hcl
terraform {
  required_providers {
    cobbler = {
      source  = "cobbler/cobbler"
      version = "~> 4.0"
    }
  }
}
```

Run `tofu init` instead of `terraform init`.
