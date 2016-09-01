package main

import (
	"flag"
)

type (
	asyncFlag   *bool
	configPath  *string
	keyPath     *string
	serverGroup *string
	command     *string
)

type input struct {
	asyncFlag   asyncFlag
	configPath  configPath
	keyPath     keyPath
	serverGroup serverGroup
	command     command
}

const (
	FLAG_ASYNC_KEY = "async"
	FLAG_ASYNC_DESC = "async"
	FLAG_ASYNC_DEF = false

	FLAG_CONFIG_PATH_KEY = "c"
	FLAG_CONFIG_PATH_DESC = "yaml config file"
	FLAG_CONFIG_PATH_DEF = "config.yml"

	FLAG_KEY_PATH_KEY = "i"
	FLAG_KEY_PATH_DESC = "identity_file"
	FLAG_KEY_PATH_DEF = ""

	FLAG_GROUP_KEY = "s"
	FLAG_GROUP_DESC = "server group"
	FLAG_GROUP_DEF = ""

	FLAG_COMMAND_KEY = "command"
	FLAG_COMMAND_DESC = "path to ssh key ie. ~/.ssh/id_rsa"
	FLAG_COMMAND_DEF = ""
)

func (this *input)parse() {
	this.asyncFlag = asyncFlag(flag.Bool(FLAG_ASYNC_KEY, FLAG_ASYNC_DEF, FLAG_ASYNC_DESC))
	this.configPath = configPath(flag.String(FLAG_CONFIG_PATH_KEY, FLAG_CONFIG_PATH_DEF, FLAG_CONFIG_PATH_DESC))
	this.keyPath = keyPath(flag.String(FLAG_KEY_PATH_KEY, FLAG_KEY_PATH_DEF, FLAG_KEY_PATH_DESC))
	this.serverGroup = serverGroup(flag.String(FLAG_GROUP_KEY, FLAG_GROUP_DEF, FLAG_GROUP_DESC))
	this.command = command(flag.String(FLAG_COMMAND_KEY, FLAG_COMMAND_DEF, FLAG_COMMAND_DESC))
	flag.Parse()
}

func (this *input)validate() {
	//this.asyncFlag.validate()
	//this.configPath.validate()
	//this.keyPath.validate()
	//this.serverGroup.validate()
	//this.command.validate()
}

//func (this *asyncFlag)validate() {
//	if *this != true || *this != false{
//		panic("asyncFlag")
//	}
//}
//
//func (this *configPath)validate() {
//	if *this == ""{
//		panic("configPath")
//	}
//}
//
//func (this *keyPath)validate() {
//	if *this == ""{
//		panic("keyPath")
//	}
//}
//
//func (this *serverGroup)validate() {
//	if *this == ""{
//		panic("serverGroup")
//	}
//}
//
//func (this *command)validate() {
//	if *this == ""{
//		panic("command")
//	}
//}
