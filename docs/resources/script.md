---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "myscribae_script Resource - myscribae"
subcategory: ""
description: |-
  
---

# myscribae_script (Resource)



## Example Usage

```terraform
resource "myscribae_script" "example" {
  provider_id        = myscribae_provider.example.id
  script_group_id    = myscribae_script_group.example.id
  alt_id             = "example_script_group"
  name               = "Example Group"
  description        = "Example group is a group of scripts"
  price_in_cents     = 1000
  sla_sec            = 3600
  token_lifetime_sec = 1800
  recurrence         = "monthly"
  public             = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `alt_id` (String) The alt id of the script
- `description` (String) The description of the script
- `name` (String) The name of the script
- `price_in_cents` (Number) The price in cents of the script (minimum 1)
- `provider_id` (String) The provider id of the script
- `recurrence` (String) The recurrence of the script
- `script_group_id` (String) The script group uuid
- `sla_sec` (Number) The SLA in seconds of the script (minimum 2400)
- `token_lifetime_sec` (Number) The token lifetime in seconds of the script (minimum 600)

### Optional

- `public` (Boolean) Is the script public

### Read-Only

- `id` (String) The id of the script
- `uuid` (String) The uuid of the script
