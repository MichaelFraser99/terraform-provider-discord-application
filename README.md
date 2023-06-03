# discord-application-terraform-provider
A terraform provider enabling the tracking of registered discord bot commands through the use of IaC

# Import Format
```shell
terraform import --var-file=vars.tfvars  discord-application_command.example "application_id-command_id"
```