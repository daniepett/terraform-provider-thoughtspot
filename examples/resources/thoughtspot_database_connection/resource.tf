resource "thoughtspot_database_connection" "this" {
  name                = "tf-test"
  description         = "example connection"
  data_warehouse_type = "SNOWFLAKE"
  data_warehouse_config = {
    configuration = {
      account_name = "account.region"
      user         = "username"
      password     = "password"
      role         = "READ_ONLY"
      warehouse    = "COMPUTE_WH"
    }
  }
  validate = false
}
