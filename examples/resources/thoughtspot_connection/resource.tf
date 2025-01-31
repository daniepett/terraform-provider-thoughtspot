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
