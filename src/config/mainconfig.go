package config

import (
	"log"

	logging "gopkg.in/tokopedia/logging.v1"
)

type Config struct {
	Server   ServerConf
	Database map[string]*struct {
		Master string
		Slave  []string
	}
	NSQ nsqStruct
}

type nsqStruct struct {
	NSQD    string
	Lookupd string
}

type ServerConf struct {
	Name string
}

type NSQ struct {
}

func InitConfig() *Config {
	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	return &cfg
}
