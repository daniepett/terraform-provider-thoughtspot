resource "thoughtspot_custom_calendar" "existing_table" {
  name           = "Custom Calendar from Existing Table"
  existing_table = true
  table_reference {
    connection_identifier = "your_connection_identifier"
    database_name         = "your_database_name"
    schema_name           = "your_schema_name"
    table_name            = "your_table_name"
  }
}

resource "thoughtspot_custom_calendar" "input_params" {
  name           = "Custom Calendar from Input Params"
  existing_table = false
  table_reference {
    connection_identifier = "your_connection_identifier"
    database_name         = "your_database_name"
    schema_name           = "your_schema_name"
    table_name            = "your_table_name"
  }
  start_date          = "09/01/2015"
  end_date            = "12/31/2030"
  calendar_type       = "FOUR_FOUR_FIVE"
  month_offset        = "September"
  start_day_of_week   = "Monday"
  quarter_name_prefix = "Q"
  year_name_prefix    = "Y"
}
