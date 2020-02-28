---
layout: "cobbler"
page_title: "Cobbler: cobbler_snippet"
sidebar_current: "docs-cobbler-resource-snippet"
description: |-
  Manages a Snippet within Cobbler.
---

# cobbler_snippet

Manages a Snippet within Cobbler.

## Example Usage

```hcl
resource "cobbler_snippet" "my_snippet" {
  name = "my_snippet"
  body = "<content of snippet>"
}
```

## Argument Reference

The following arguments are supported:

* `body` - (Required) The body of the snippet.  May also point
  to a file: `body = file("my_snippet")`

* `name` - (Required) The name of the snippet. This must be the name
  only, so without `/var/lib/cobbler/snippets`.
