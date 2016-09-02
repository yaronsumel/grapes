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

const (
	FLAG_ASYNC_KEY = "async"
	FLAG_ASYNC_DESC = "async - when true, parallel executing over servers"
	FLAG_ASYNC_DEF = false

	FLAG_CONFIG_PATH_KEY = "c"
	FLAG_CONFIG_PATH_DESC = "config file - yaml config file"
	FLAG_CONFIG_PATH_DEF = ""

	FLAG_KEY_PATH_KEY = "i"
	FLAG_KEY_PATH_DESC = "identity file - path to private key"
	FLAG_KEY_PATH_DEF = ""

	FLAG_GROUP_KEY = "s"
	FLAG_GROUP_DESC = "server group - name of the server group"
	FLAG_GROUP_DEF = ""

	FLAG_COMMAND_KEY = "cmd"
	FLAG_COMMAND_DESC = "command name - name of the command to run"
	FLAG_COMMAND_DEF = ""
)

func (this *input)parse() {

	asyncFlagPtr := flag.Bool(FLAG_ASYNC_KEY, FLAG_ASYNC_DEF, FLAG_ASYNC_DESC)
	configPathPtr := flag.String(FLAG_CONFIG_PATH_KEY, FLAG_CONFIG_PATH_DEF, FLAG_CONFIG_PATH_DESC)
	keyPathPtr := flag.String(FLAG_KEY_PATH_KEY, FLAG_KEY_PATH_DEF, FLAG_KEY_PATH_DESC)
	serverGroupPtr := flag.String(FLAG_GROUP_KEY, FLAG_GROUP_DEF, FLAG_GROUP_DESC)
	commandPtr := flag.String(FLAG_COMMAND_KEY, FLAG_COMMAND_DEF, FLAG_COMMAND_DESC)

	flag.Parse()

	this.asyncFlag = asyncFlag(*asyncFlagPtr)
	this.commandName = commandName(*commandPtr)
	this.serverGroup = serverGroup(*serverGroupPtr)
	this.keyPath = keyPath(*keyPathPtr)
	this.configPath = configPath(*configPathPtr)

}

func (this *input)validate() {
	this.configPath.validate()
	this.keyPath.validate()
	this.serverGroup.validate()
	this.commandName.validate()
}

func (val *configPath)validate() {
	if *val == "" {
		fatal("configPath is empty please set grapes -c config.yml")
	}
}

func (val *keyPath)validate() {
	if *val == "" {
		fatal("idendity file path is empty please set grapes -i ~/.ssh/id_rsaa")
	}
}

func (val *serverGroup)validate() {
	if *val == "" {
		fatal("server group is empty please set grapes -s server_group")
	}
}

func (val *commandName)validate() {
	if *val == "" {
		fatal("command name is empty please set grapes -command whats_up")
	}
}
