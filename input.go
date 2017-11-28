package main

import (
	"errors"
	"flag"

	"github.com/mitchellh/go-homedir"
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

func getInputData() *input {

	verifyActionFlagPtr := flag.Bool("y", false, "force yes")
	asyncFlagPtr := flag.Bool("async", false, "async - if true, parallel executing over servers")
	configPathPtr := flag.String("c", "", "config file - yaml config file, defaulting to ~/.grapes.yml")
	keyPathPtr := flag.String("i", "", "identity file - path to private key, defaulting to ~/.ssh/id_rsa")
	serverGroupPtr := flag.String("s", "", "server group - name of the server group")
	commandPtr := flag.String("cmd", "", "command name - name of the command to run")

	flag.Parse()

	if *configPathPtr == "" {
		potentialConfig, err := homedir.Expand("~/.grapes.yml")
		if err == nil {
			configPathPtr = &potentialConfig
		}
	}

	if *keyPathPtr == "" {
		potentialKeyPath, err := homedir.Expand("~/.ssh/id_rsa")
		if err == nil {
			keyPathPtr = &potentialKeyPath
		}
	}

	return &input{
		verifyFlag:  verifyFlag(*verifyActionFlagPtr),
		asyncFlag:   asyncFlag(*asyncFlagPtr),
		commandName: commandName(*commandPtr),
		serverGroup: serverGroup(*serverGroupPtr),
		keyPath:     keyPath(*keyPathPtr),
		configPath:  configPath(*configPathPtr),
	}
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
		return input.newError("identity file path is empty please set grapes -i ~/.ssh/id_rsa")
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
