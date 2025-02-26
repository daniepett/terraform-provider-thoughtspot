resource "thoughtspot_share_metadata" "this" {
  metadata_type        = "LIVEBOARD"
  metadata_identifiers  = ["Test 1", "Test 2"]
  principal_type       = "USER_GROUP"
  principal_identifier = "TEST"
  discoverable         = true
  share_mode           = "READ_ONLY"
}