package main

import (
	"errors"
	"fmt"
	"io"
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

func (app *grape) run(knownHostsPath string, out io.Writer) error {

	knownHostsPath, err := homedir.Expand(knownHostsPath)
	if err != nil {
		return err
	}

	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var outputMu sync.Mutex
	errs := make([]error, len(app.servers))
	runServer := func(i int, server server) {
		result, runErr := app.runOnServer(server, hostKeyCallback)
		outputMu.Lock()
		outputErr := result.printOutput(out)
		outputMu.Unlock()
		errs[i] = errors.Join(runErr, outputErr)
	}
	for i, server := range app.servers {
		if app.input.asyncFlag {
			wg.Add(1)
			go func() {
				defer wg.Done()
				runServer(i, server)
			}()
		} else {
			runServer(i, server)
		}
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (app *grape) runOnServer(server server, hostKeyCallback ssh.HostKeyCallback) (server, error) {
	client, err := app.ssh.newClient(server, hostKeyCallback)
	if err != nil {
		server.Fatal = err.Error()
		return server, fmt.Errorf("%s: %w", server.Host, err)
	}
	defer client.Close()

	server.Output = client.execCommands(app.commands)
	var errs []error
	for _, output := range server.Output {
		if output.Error != "" {
			errs = append(errs, fmt.Errorf("%s: command %q: %s", server.Host, output.Command, output.Error))
		}
	}
	return server, errors.Join(errs...)
}

func (s *server) printOutput(out io.Writer) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, string(data)); err != nil {
		return fmt.Errorf("%s: could not print results: %w", s.Host, err)
	}
	return nil
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
