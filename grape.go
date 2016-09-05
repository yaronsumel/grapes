package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type (
	grape struct {
		input  input
		ssh    grapeSSH
		config config
	}
	serverOutput struct {
		Server server             `yaml:"server"`
		Output []*grapeCommandStd `yaml:"stds"`
	}
	serverFatal struct {
		Server server `yaml:"server"`
		Fatal  string `yaml:"fatal"`
	}
)

var wg sync.WaitGroup

func (app *grape) init() {
	app.input.parse()
	if err := app.input.validate(); err != nil {
		panic(err)
	}
	if err := app.config.set(app.input.configPath); err != nil {
		panic(err)
	}
	if err := app.ssh.setKey(app.input.keyPath); err != nil {
		panic(err)
	}
}

func (app *grape) run() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\r\nFatal: %s \n\n", err)
			os.Exit(1)
		}
	}()
	app.init()
	servers, err := app.config.getServersFromConfig(app.input.serverGroup)
	if err != nil {
		panic(err)
	}
	app.input.verifyAction(servers)
	for _, server := range servers {
		if app.input.asyncFlag {
			wg.Add(1)
			go app.runCommandAsync(server, &wg)
		} else {
			app.runCommand(server)
		}
	}
	if app.input.asyncFlag {
		wg.Wait()
	}
}

func (app *grape) runCommandAsync(server server, wg *sync.WaitGroup) {
	a, err := app.runCommandsOnServer(server)
	wg.Done()
	if err != nil {
		app.print(serverFatal{
			Server: server,
			Fatal:  err.Error(),
		})
		return
	}
	app.print(a)
}

func (app *grape) runCommand(server server) {
	a, err := app.runCommandsOnServer(server)
	if err != nil {
		app.print(serverFatal{
			Server: server,
			Fatal:  err.Error(),
		})
		return
	}
	app.print(a)
}

func (app *grape) runCommandsOnServer(server server) (*serverOutput, error) {
	commands, err := app.config.getCommandsFromConfig(app.input.commandName)
	if err != nil {
		return nil, err
	}
	client, err := app.ssh.newClient(server)
	if err != nil {
		return nil, err
	}
	so := serverOutput{
		Server: server,
	}
	for _, command := range commands {
		grapeCommandStdOut, err := client.exec(command)
		if err != nil {
		}
		so.Output = append(so.Output, grapeCommandStdOut)
	}
	return &so, nil
}

func (app grape) print(a interface{}) {
	out, err := yaml.Marshal(a)
	if err != nil {
		panic("something went wrong with the output")
	}
	fmt.Println(string(out))
}

func (so *serverOutput) print() {
	out, err := yaml.Marshal(so)
	if err != nil {
		panic("something went wrong with the output")
	}
	fmt.Println(string(out))
}
