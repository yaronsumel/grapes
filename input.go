package main

import (
	"flag"
)

type (
	asyncFlag   bool
	configPath  string
	keyPath     string
	serverGroup string
	commandName string
)

type input struct {
	asyncFlag   asyncFlag
	configPath  configPath
	keyPath     keyPath
	serverGroup serverGroup
	commandName commandName
}

func (input *input) parse() {

	asyncFlagPtr := flag.Bool("async", false, "async - when true, parallel executing over servers")
	configPathPtr := flag.String("c", "", "config file - yaml config file")
	keyPathPtr := flag.String("i", "", "identity file - path to private key")
	serverGroupPtr := flag.String("s", "", "server group - name of the server group")
	commandPtr := flag.String("cmd", "", "command name - name of the command to run")

	flag.Parse()

	input.asyncFlag = asyncFlag(*asyncFlagPtr)
	input.commandName = commandName(*commandPtr)
	input.serverGroup = serverGroup(*serverGroupPtr)
	input.keyPath = keyPath(*keyPathPtr)
	input.configPath = configPath(*configPathPtr)

}

func (input *input) validate() {
	input.configPath.validate()
	input.keyPath.validate()
	input.serverGroup.validate()
	input.commandName.validate()
}

func (val *configPath) validate() {
	if *val == "" {
		fatal("configPath is empty please set grapes -c config.yml")
	}
}

func (val *keyPath) validate() {
	if *val == "" {
		fatal("idendity file path is empty please set grapes -i ~/.ssh/id_rsaa")
	}
}

func (val *serverGroup) validate() {
	if *val == "" {
		fatal("server group is empty please set grapes -s server_group")
	}
}

func (val *commandName) validate() {
	if *val == "" {
		fatal("command name is empty please set grapes -cmd whats_up")
	}
}
