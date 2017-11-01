# Check Consul Service [![Build Status](https://travis-ci.org/ujenmr/check-consul-service.svg?branch=master)](https://travis-ci.org/ujenmr/check-consul-service)

Nagios/Icinga plugin checks consul alive services

## Usage:
```bash
./check_consul_service -consul-addr 127.0.0.1:8500 -w 1 -c 0
```

## Icinga configuration

commands.conf:
```
object CheckCommand "consul-service" {
  command = [ SysconfDir + "/icinga2/scripts/check_consul_service" ]

  arguments = {
    "-consul-addr" = {
      required = true
      value = "$consul_address$"
    }
    "-user" = {
      value = "$consul_auth_user$"
      description = "Consul Auth User"
    }
    "-password" = {
      value = "$consul_auth_password$"
      description = "Consul Auth Password"
    }
    "-scheme" = {
      value = "$consul_scheme$"
      description = "Consul Scheme (http/https)"
    }
    "-w" = {
      value = "$warning$"
    }
    "-c" = {
      value = "$critical$"
    }
  }

  vars.consul_address = "$address$:8500"
  vars.warning = "0"
  vars.critical = "0"
}
```
