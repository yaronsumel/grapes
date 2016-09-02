package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	Servers  []server
	Commands []command
	command  string
	server   struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		User string `yaml:"user"`
	}
	config struct {
		version  string              `yaml:"version"`
		Servers  map[string]Servers  `yaml:"servers"`
		Commands map[string]Commands `yaml:"commands"`
	}
)

func (c *config) set(configPath configPath) {
	data, err := ioutil.ReadFile(string(configPath))
	if err != nil {
		fatal(fmt.Sprintf("Could not open %s  ", configPath))
	}
	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		fatal(fmt.Sprintf("Could not parse config file. make sure its yaml."))
	}
}

func (config *config) getServersFromConfig(serverGroup serverGroup) Servers {
	group, ok := config.Servers[string(serverGroup)]
	if !ok {
		fatal(fmt.Sprintf("Could not find [%s] in server group.", serverGroup))
	}
	return group
}

func (config *config) getCommandsFromConfig(commandName commandName) Commands {
	commands, ok := config.Commands[string(commandName)]
	if !ok {
		fatal(fmt.Sprintf("Command %s was not found.", commandName))
	}
	return commands
}
