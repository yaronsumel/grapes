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

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		cOk := config{}
		cOk.set(configPath("testFiles/demo_config.yml"))
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		cNotOk := config{}
		cNotOk.set(configPath("not_demo_config.yml"))
		t.FailNow()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		cNotOk1 := config{}
		cNotOk1.set(configPath("main.go"))
		t.FailNow()
	}()

}

func TestGetServersFromConfig(t *testing.T) {
	groupName := "test1"
	serversArray := demoConfig.getServersFromConfig(serverGroup(groupName))
	if !reflect.DeepEqual(demoConfig.Servers["test1"], serversArray) {
		t.FailNow()
	}

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoConfig.getServersFromConfig(serverGroup("false_group_name"))
		t.FailNow()
	}()

}

func TestGetCommandsFromConfig(t *testing.T) {
	cmdName := "test_cmd"
	cmds := demoConfig.getCommandsFromConfig(commandName(cmdName))
	if !reflect.DeepEqual(demoConfig.Commands["test_cmd"], cmds) {
		t.FailNow()
	}
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoConfig.getCommandsFromConfig(commandName("false_command"))
		t.FailNow()
	}()
}
