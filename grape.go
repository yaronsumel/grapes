package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"gopkg.in/yaml.v2"
)

type (
	grape struct {
		input    input
		ssh      grapeSSH
		config   config
		servers  servers
		commands commands
	}
)

var wg sync.WaitGroup

func newGrape(input *input) *grape {
	app := grape{}
	var err error

	app.input = *input

	//validate input
	if err = app.input.validate(); err != nil {
		panic(err)
	}
	//set config into place
	if err = app.config.set(app.input.configPath); err != nil {
		panic(err)
	}
	// data !
	if app.servers, err = app.config.getServersFromConfig(app.input.serverGroup); err != nil {
		panic(err)
	}
	if app.commands, err = app.config.getCommandsFromConfig(app.input.commandName); err != nil {
		panic(err)
	}
	//load private key
	if err = app.ssh.setKey(app.input.keyPath); err != nil {
		panic(err)
	}
	return &app
}

func (app *grape) run(knownHostsPath string) {

	knownHostsPath, err := homedir.Expand(knownHostsPath)
	if err != nil {
		log.Fatal(err)
	}

	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range app.servers {
		if app.input.asyncFlag {
			wg.Add(1)
			go app.runOnServer(server, hostKeyCallback, &wg)
		} else {
			app.runOnServer(server, hostKeyCallback, nil)
		}
	}
	if app.input.asyncFlag {
		wg.Wait()
	}
}

func (app *grape) runOnServer(server server, hostKeyCallback ssh.HostKeyCallback, wg *sync.WaitGroup) {

	client, err := app.ssh.newClient(server, hostKeyCallback)
	if err != nil {
		server.Fatal = err.Error()
	} else {
		server.Output = client.execCommands(app.commands)
	}
	if wg != nil {
		wg.Done()
	}
	server.printOutput()
}

func (s *server) printOutput() {
	out, _ := yaml.Marshal(s)
	fmt.Println(string(out))
}

func (app *grape) verifyAction() {
	var char = "n"
	fmt.Println("The following command will run on the following servers:")
	fmt.Printf("command `%s` will run over `%s`.\n", app.input.commandName, app.input.serverGroup)
	fmt.Println("commands:")
	for k, v := range app.commands {
		fmt.Printf("\t#%d - `%s` \n", k, v)
	}
	fmt.Println("servers:")
	for k, v := range app.servers {
		fmt.Printf("\t#%d - %s [%s@%s] \n", k, v.Name, v.User, v.Host)
	}
	if app.input.verifyFlag {
		fmt.Println("-y used.forced to continue.")
		return
	}
	fmt.Print("\n -- are your sure? [y/N] : ")
	if _, err := fmt.Scanf("%s", &char); err != nil || char != "y" {
		panic("type y to continue")
	}
}
