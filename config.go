package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	servers  []server
	commands []command
	command  string
	server   struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		User string `yaml:"user"`
	}
	config struct {
		Version  string              `yaml:"version"`
		Servers  map[string]servers  `yaml:"servers"`
		Commands map[string]commands `yaml:"commands"`
	}
)

func (conf *config) set(configPath configPath) {
	data, err := ioutil.ReadFile(string(configPath))
	if err != nil {
		panic(fmt.Sprintf("Could not open %s  ", configPath))
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		panic(fmt.Sprintf("Could not parse config file. make sure its yaml."))
	}
}

func (conf *config) getServersFromConfig(serverGroup serverGroup) servers {
	group, ok := conf.Servers[string(serverGroup)]
	if !ok {
		panic(fmt.Sprintf("Could not find [%s] in server group.", serverGroup))
	}
	return group
}

func (conf *config) getCommandsFromConfig(commandName commandName) commands {
	commands, ok := conf.Commands[string(commandName)]
	if !ok {
		panic(fmt.Sprintf("Command %s was not found.", commandName))
	}
	return commands
}