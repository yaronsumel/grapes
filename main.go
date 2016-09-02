package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"sync"
)

type (
	grape struct {
		input  input
		ssh    grapeSsh
		config config
	}
	serverOutput struct {
		Server server             `yaml:server`
		Output []*grapeCommandStd `yaml:stds`
	}
)

const (
	WELCOME = `
//      ____ __________ _____  ___  _____
//     / __  / ___/ __  / __ \/ _ \/ ___/
//    / /_/ / /  / /_/ / /_/ /  __(__  )
//    \__, /_/   \__,_/ .___/\___/____/
//   /____/          /_/ v 0.1.2 // Yaron Sumel [yaronsu@gmail.com]
//
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
	grape.input.validate()

	grape.ssh.setKey(grape.input.keyPath)
	grape.config.set(grape.input.configPath)

	return &grape
}

func (app *grape) runApp() {
	servers := app.config.getServersFromConfig(app.input.serverGroup)
	//todo verify action
	for _, server := range servers {
		if app.input.asyncFlag {
			wg.Add(1)
			go app.asyncRunCommand(server, &wg)
		} else {
			app.runCommandsOnServer(server)
		}
	}
	if app.input.asyncFlag {
		wg.Wait()
	}
}

func (app *grape) asyncRunCommand(server server, wg *sync.WaitGroup) {
	app.runCommandsOnServer(server)
	wg.Done()
}

func (app *grape) runCommandsOnServer(server server) {

	commands := app.config.getCommandsFromConfig(app.input.commandName)
	client := app.ssh.newClient(server)
	so := serverOutput{
		Server: server,
	}

	for _, command := range commands {
		so.Output = append(so.Output, client.newSession().exec(command))
	}

	//done with all commands for this server
	so.print()
}

func (so *serverOutput) print() {
	out, err := yaml.Marshal(so)
	if err != nil {
		fatal("something went wrong with the output")
	}
	fmt.Println(string(out))
}
