package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
)

type (
	Servers []server
	Commands []command
	server struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		User string `yaml:"user"`
	}
	config struct {
		version  string `yaml:"version"`
		Servers  map[string]Servers `yaml:"servers"`
		Commands map[string]Commands `yaml:"commands"`
	}
)

func (c *config)set(configPath configPath) {
	data, err := ioutil.ReadFile(string(*configPath))
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (config *config)getServersFromConfig(serverGroup serverGroup) Servers {
	group, ok := config.Servers[string(*serverGroup)]
	if !ok {
		panic("cant find that server group")
	}
	return group
}

func (config *config)getCommandsFromConfig(commandName command) Commands {
	commands, ok := config.Commands[string(*commandName)]
	if !ok {
		panic("cant find that command")
	}
	return commands
}