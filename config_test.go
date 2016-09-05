package main

import (
	"reflect"
	"testing"
)

func getDemoConfig() config {
	config := config{
		Servers:  make(map[string]servers),
		Commands: make(map[string]commands),
	}
	config.Servers["test1"] = []server{
		{
			Host: "test1",
			Name: "test1",
			User: "test1",
		},
		{
			Host: "test12",
			Name: "test12",
			User: "test12",
		},
	}
	config.Commands["test_cmd"] = []command{
		command("ls -al #1"),
		command("ls -al #2"),
	}
	return config
}

var demoConfig = getDemoConfig()

func TestSet(t *testing.T) {
	cOk := config{}
	if cOk.set(configPath("testFiles/demo_config.yml")) != nil {
		t.FailNow()
	}
	cNotOk := config{}
	if cNotOk.set(configPath("not_demo_config.yml")) == nil {
		t.FailNow()
	}
	cNotOk1 := config{}
	if cNotOk1.set(configPath("main.go")) == nil {
		t.FailNow()
	}
}

func TestGetServersFromConfig(t *testing.T) {
	serversArray, err := demoConfig.getServersFromConfig(serverGroup("test1"))
	if err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(demoConfig.Servers["test1"], serversArray) {
		t.FailNow()
	}
	if _, err = demoConfig.getServersFromConfig(serverGroup("false_group_name")); err == nil {
		t.FailNow()
	}
}

func TestGetCommandsFromConfig(t *testing.T) {
	cmds, err := demoConfig.getCommandsFromConfig(commandName("test_cmd"))
	if err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(demoConfig.Commands["test_cmd"], cmds) {
		t.FailNow()
	}
	if _, err := demoConfig.getCommandsFromConfig(commandName("false_command")); err == nil {
		t.FailNow()
	}
}
