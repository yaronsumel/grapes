package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	servers  []server
	commands []command
	command  string
	server   struct {
		Name   string         `yaml:"name"`
		Host   string         `yaml:"host"`
		User   string         `yaml:"user"`
		Fatal  string         `yaml:"fatal"`
		Output sshOutputArray `yaml:"output"`
	}
	config struct {
		Version  string              `yaml:"version"`
		Servers  map[string]servers  `yaml:"servers"`
		Commands map[string]commands `yaml:"commands"`
	}
	configError error
)

func (conf *config) newError(errMsg string) configError {
	return errors.New(errMsg)
}

func (conf *config) set(configPath configPath) error {
	data, err := ioutil.ReadFile(string(configPath))
	if err != nil {
		return conf.newError(fmt.Sprintf("Could not open %s  ", configPath))
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		return conf.newError(fmt.Sprintf("Could not parse config file. Make sure its yaml."))
	}
	return nil
}

func (conf *config) getServersFromConfig(serverGroup serverGroup) (servers, configError) {
	group, ok := conf.Servers[string(serverGroup)]
	if !ok {
		return nil, conf.newError(fmt.Sprintf("Could not find [%s] in server group.", serverGroup))
	}
	return group, nil
}

func (conf *config) getCommandsFromConfig(commandName commandName) (commands, configError) {
	commands, ok := conf.Commands[string(commandName)]
	if !ok {
		return nil, conf.newError(fmt.Sprintf("Command %s was not found.", commandName))
	}
	return commands, nil
}
