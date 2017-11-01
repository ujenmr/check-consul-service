// Copyright 2016 Ievgen Khmelenko <ujenmr@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
  Usage of ./check_consul_service:
  -c int
    	Critical
  -consul-addr string
    	Consul Address (default "127.0.0.1:8500")
  -password string
    	Consul Auth Password
  -scheme string
    	Consul Scheme (default "http")
  -user string
    	Consul Auth User
  -w int
    	Warning (default 1)
*/
package main

import (
	"flag"
	"fmt"
	"os"

	api "github.com/hashicorp/consul/api"
)

const CODE_OK = 0
const CODE_WARNING = 1
const CODE_CRITICAL = 2
const CODE_UNKNOWN = 3

var consulAddr string
var consulUser string
var consulPass string
var consulScheme string
var warningLimit int
var criticalLimit int

func printNagiosOut(msg string, code int) {

	/*

		According by https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html

		,-------------------------------------------------------------,
		| Plugin Return Code | Service State | Host State             |
		|--------------------|---------------|------------------------|
		| 0	                 | OK	         | UP                     |
		| 1	                 | WARNING	     | UP or DOWN/UNREACHABLE |
		| 2	                 | CRITICAL	     | DOWN/UNREACHABLE       |
		| 3	                 | UNKNOWN	     | DOWN/UNREACHABLE       |
		'-------------------------------------------------------------'

	*/

	if code == CODE_OK {
		fmt.Printf("CONSUL-SERVICE OK: %s", msg)
	}

	if code == CODE_WARNING {
		fmt.Printf("CONSUL-SERVICE WARNING: %s", msg)
	}

	if code == CODE_CRITICAL {
		fmt.Printf("CONSUL-SERVICE CRITICAL: %s", msg)
	}

	if code == CODE_UNKNOWN {
		fmt.Printf("CONSUL-SERVICE UNKNOWN: %s", msg)
	}

	os.Exit(code)
}

func init() {
	flag.StringVar(&consulAddr, "consul-addr", "127.0.0.1:8500", "Consul Address")
	flag.StringVar(&consulUser, "user", "", "Consul Auth User")
	flag.StringVar(&consulPass, "password", "", "Consul Auth Password")
	flag.StringVar(&consulScheme, "scheme", "http", "Consul Scheme")
	flag.IntVar(&warningLimit, "w", 1, "Warning")
	flag.IntVar(&criticalLimit, "c", 0, "Critical")
	flag.Parse()
	if warningLimit < criticalLimit {
		printNagiosOut("Warning value must be less then critical", CODE_UNKNOWN)
	}
}

func main() {
	config := api.DefaultConfig()
	config.Address = consulAddr
	if consulUser != "" && consulPass != "" {
		config.HttpAuth = &api.HttpBasicAuth{Username: consulUser, Password: consulPass}
	}
	if consulScheme != "http" {
		config.Scheme = consulScheme
	}
	client, err := api.NewClient(config)
	if err != nil {
		printNagiosOut(err.Error(), CODE_UNKNOWN)
	}

	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)

	if err != nil {
		printNagiosOut(err.Error(), CODE_UNKNOWN)
	}

	var s string
	exitCode := CODE_OK
	for name, _ := range services {
		svcDesc, _, err := catalog.Service(name, "", nil)
		if err != nil {
			printNagiosOut(err.Error(), CODE_UNKNOWN)
		}
		if len(svcDesc) <= warningLimit {
			exitCode = CODE_WARNING
			if len(svcDesc) <= criticalLimit {
				exitCode = CODE_CRITICAL
			}
		}

		s += fmt.Sprintf("%s=%d ", name, len(svcDesc))
	}
	printNagiosOut(s, exitCode)
}
