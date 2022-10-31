package main

import (
	apiconfig "github.com/Alex-Eftimie/api-config"
)

type ConfigType struct {
	apiconfig.Configuration
	ProxyAddr string
	BindPort  int
}

var Co *ConfigType

func init() {
	// log.Println("Reading Config")
	Co = &ConfigType{
		Configuration: *apiconfig.NewConfig("config.jsonc"),
	}

	apiconfig.LoadConfig(Co)

}
