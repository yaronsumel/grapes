package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	testPath := configPath("demo_config.yml")
	c1 := config{}
	c2 := config{}
	c1.set(testPath)
	data, err := ioutil.ReadFile(string(testPath))
	if err != nil {
		t.FailNow()
	}
	if err := yaml.Unmarshal([]byte(data), &c2); err != nil {
		t.FailNow()
	}
	if !reflect.DeepEqual(c1, c2) {
		t.FailNow()
	}
}

func TestGetServersFromConfig(t *testing.T) {
	groupName := "test1"
	serversArray := demoConfig.getServersFromConfig(serverGroup(groupName))
	if !reflect.DeepEqual(demoConfig.Servers["test1"], serversArray) {
		t.FailNow()
	}
}

func TestGetCommandsFromConfig(t *testing.T) {
	cmdName := "test_cmd"
	cmds := demoConfig.getCommandsFromConfig(commandName(cmdName))
	if !reflect.DeepEqual(demoConfig.Commands["test_cmd"], cmds) {
		t.FailNow()
	}
}
