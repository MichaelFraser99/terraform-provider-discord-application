resource "discord-application_command" "poke" {
  application_id = "9876543210123456789"
  name = "poke"
  description = "a poke to the application"
  type = 1
}