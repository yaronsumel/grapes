package main

//     ____ __________ _____  ___
//    / __ // ___/ __ // __ \/ _ \
//   / /_/ / /  / /_/ / /_/ /  __/
//   \__, /_/   \__,_/ .___/\___/
//  /____/          /_/

import (
	"sync"
	"fmt"
	"gopkg.in/yaml.v2"
)

type (
	grape struct {
		input  input
		auth   auth
		config config
	}
	serverOutput struct {
		Server server `yaml:server`
		Output []*grapeCommandStd `yaml:stds`
	}
)

const (
	WELCOME = `
//
//     ____ __________ _____  ___
//    / __ // ___/ __ // __ \/ _ \
//   / /_/ / /  / /_/ / /_/ /  __/
//   \__, /_/   \__,_/ .___/\___/
//  /____/          /_/
//
//
// Version 0.1 -- Yaron Sumel [yaronsu@gmail.com]
`
)

var wg sync.WaitGroup

func main() {
	fmt.Println(WELCOME)
	newGrape().runApp()
}

func newGrape() *grape {

	grape := grape{}
	//parse flags and validate it
	grape.input.parse()
	//grape.input.validate()

	grape.auth.setKey(grape.input.keyPath)
	grape.config.set(grape.input.configPath)

	return &grape
}

func (app *grape)runApp() {
	servers := app.config.getServersFromConfig(app.input.serverGroup)
	//todo verify action
	for _, server := range servers {
		if *app.input.asyncFlag {
			wg.Add(1)
			go app.asyncRunCommand(server, &wg)
		} else {
			app.runCommandsOnServer(server)
		}
	}
	if *app.input.asyncFlag {
		wg.Wait()
	}
}

func (app *grape)asyncRunCommand(server server, wg *sync.WaitGroup) {
	app.runCommandsOnServer(server)
	wg.Done()
}

func (app *grape)runCommandsOnServer(server server) {
	client := app.auth.newClient(server)
	o := serverOutput{
		Server:server,
	}
	for _, command := range app.config.getCommandsFromConfig(app.input.command) {
		o.Output = append(o.Output, client.newSession().exec(command))
	}
	//done with all commands for this server
	out, _ := yaml.Marshal(o)
	fmt.Println(string(out))
}