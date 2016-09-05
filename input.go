package main

import (
	"errors"
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
	inputError error
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

func (input *input) newError(errMsg string) inputError {
	return errors.New(errMsg)
}

func (input *input) validate() inputError {
	if err := input.configPath.validate(input); err != nil {
		return err
	}
	if err := input.keyPath.validate(input); err != nil {
		return err
	}
	if err := input.serverGroup.validate(input); err != nil {
		return err
	}
	if err := input.commandName.validate(input); err != nil {
		return err
	}
	return nil
}

func (val *configPath) validate(input *input) inputError {
	if *val == "" {
		return input.newError("configPath is empty please set grapes -c config.yml")
	}
	return nil
}

func (val *keyPath) validate(input *input) inputError {
	if *val == "" {
		return input.newError("idendity file path is empty please set grapes -i ~/.ssh/id_rsaa")
	}
	return nil
}

func (val *serverGroup) validate(input *input) inputError {
	if *val == "" {
		return input.newError("server group is empty please set grapes -s server_group")
	}
	return nil
}

func (val *commandName) validate(input *input) inputError {
	if *val == "" {
		return input.newError("command name is empty please set grapes -cmd whats_up")
	}
	return nil
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
