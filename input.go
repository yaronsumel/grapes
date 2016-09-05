package main

import (
	"flag"
	"fmt"
)

type (
	asyncFlag   bool
	configPath  string
	keyPath     string
	serverGroup string
	commandName string
	verifyFlag  bool

	input struct {
		asyncFlag   asyncFlag
		configPath  configPath
		keyPath     keyPath
		serverGroup serverGroup
		commandName commandName
		verifyFlag  verifyFlag
	}
)

func (input *input) parse() {

	verifyActionFlagPtr := flag.Bool("y", false, "force yes")
	asyncFlagPtr := flag.Bool("async", false, "async - when true, parallel executing over servers")
	configPathPtr := flag.String("c", "", "config file - yaml config file")
	keyPathPtr := flag.String("i", "", "identity file - path to private key")
	serverGroupPtr := flag.String("s", "", "server group - name of the server group")
	commandPtr := flag.String("cmd", "", "command name - name of the command to run")

	flag.Parse()

	input.verifyFlag = verifyFlag(*verifyActionFlagPtr)
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
		panic("configPath is empty please set grapes -c config.yml")
	}
}

func (val *keyPath) validate() {
	if *val == "" {
		panic("idendity file path is empty please set grapes -i ~/.ssh/id_rsaa")
	}
}

func (val *serverGroup) validate() {
	if *val == "" {
		panic("server group is empty please set grapes -s server_group")
	}
}

func (val *commandName) validate() {
	if *val == "" {
		panic("command name is empty please set grapes -cmd whats_up")
	}
}

func (input *input) verifyAction(servers servers) {

	var char = "n"
	fmt.Printf("command [%s] will run on the following servers:\n", input.commandName)

	for k, v := range servers {
		fmt.Printf("\t#%d - %s [%s@%s] \n", k, v.Name, v.User, v.Host)
	}

	if input.verifyFlag {
		fmt.Println("-y used.forced to continue.")
		return
	}

	fmt.Print("\n -- are your sure? [y/N] : ")
	if _, err := fmt.Scanf("%s", &char); err != nil || char != "y" {
		panic("type y to continue")
	}
}
