---
layout: "cobbler"
page_title: "Cobbler: cobbler_template_file"
sidebar_current: "docs-cobbler-resource-template_file"
description: |-
  Manages a template File within Cobbler.
---

# cobbler_template_file

Manages a Template File within Cobbler.

## Example Usage

```hcl
resource "cobbler_template_file" "my_template" {
  name = "my_template.ks"
  body = "<content of template file>"
}
```

## Argument Reference

The following arguments are supported:

* `body` - (Required) The body of the template file.

* `name` - (Required) The name of the template file. This must be
  the name only, so without `/var/lib/cobbler/templates`.
