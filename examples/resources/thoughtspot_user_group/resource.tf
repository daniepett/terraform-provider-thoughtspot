resource "thoughtspot_user_group" "this" {
  name         = "Test Group 2"
  display_name = "Not my group"
  description  = "This is a group"
  visibility   = "SHARABLE"
  type         = "LOCAL_GROUP"
  privileges   = ["AUTHORING", "SHAREWITHALL"]
}
