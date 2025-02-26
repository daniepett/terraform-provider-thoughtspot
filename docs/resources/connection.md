---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "thoughtspot_connection Resource - terraform-provider-thoughtspot"
subcategory: ""
description: |-
  
---

# thoughtspot_connection (Resource)



## Example Usage

```terraform
resource "thoughtspot_connection" "this" {
  name        = "tf-test"
  description = "example connection"
  snowflake {
    account_name = "account.region"
    user         = "username"
    password     = "password"
    role         = "READ_ONLY"
    warehouse    = "COMPUTE_WH"
  }
  validate = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `validate` (Boolean)

### Optional

- `description` (String)
- `external_databases` (Attributes List) (see [below for nested schema](#nestedatt--external_databases))
- `redshift` (Block, Optional) (see [below for nested schema](#nestedblock--redshift))
- `snowflake` (Block, Optional) (see [below for nested schema](#nestedblock--snowflake))

### Read-Only

- `data_warehouse_type` (String)
- `id` (String) The ID of this resource.

<a id="nestedatt--external_databases"></a>
### Nested Schema for `external_databases`

Required:

- `name` (String)


<a id="nestedblock--redshift"></a>
### Nested Schema for `redshift`

Optional:

- `account_name` (String)


<a id="nestedblock--snowflake"></a>
### Nested Schema for `snowflake`

Required:

- `authentication_type` (String)

Optional:

- `access_token_url` (String)
- `account_name` (String)
- `auth_url` (String)
- `database` (String)
- `oauth_client_id` (String)
- `oauth_client_secret` (String, Sensitive)
- `passphrase` (String, Sensitive) Passphrase for the Private Key
- `password` (String, Sensitive)
- `private_key` (String, Sensitive) Private Key in PKCS8 Format
- `role` (String)
- `scope` (String)
- `user` (String)
- `warehouse` (String)
