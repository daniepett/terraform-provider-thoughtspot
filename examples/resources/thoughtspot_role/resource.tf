resource "thoughtspot_role" "this" {
  name        = "Role 1"
  description = "This is a role"
  privileges  = ["AUTHORING", "SHAREWITHALL"]
}
