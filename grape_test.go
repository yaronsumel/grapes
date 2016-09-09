package main

import (
	"testing"
)

func TestNewGrape(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "test",
			configPath:  "test",
			keyPath:     "test",
			serverGroup: "test",
		})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "staging",
		})
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "not_working",
		})
		t.FailNow()
	}()
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "not_working",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa",
			serverGroup: "staging",
		})
		t.FailNow()
	}()

	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		newGrape(&input{
			verifyFlag:  true,
			asyncFlag:   true,
			commandName: "date",
			configPath:  "testFiles/demo_config.yml",
			keyPath:     "testFiles/id_rsa.pub",
			serverGroup: "staging",
		})
		t.FailNow()
	}()

}

func TestRun(t *testing.T) {
	app := grape{}
	app.servers = servers{
		server{
			Host: "fail",
		},
		server{
			Host: "sdf.org:22",
			User: "new",
			Name: "public",
		},
	}
	app.input.asyncFlag = true
	app.run()
	app.input.asyncFlag = false
	app.run()
}

func TestPrintOutput(t *testing.T) {
	o := sshOutput{
		Command: "a",
		Std: std{
			Err: "err",
			Out: "out",
		},
	}
	s := server{
		User:  "a",
		Name:  "b",
		Host:  "c",
		Fatal: "d",
		Output: sshOutputArray{
			&o,
		},
	}
	s.printOutput()
}

func TestVerifyAction(t *testing.T) {

	app := grape{
		commands: []command{
			command("ls -al #1"),
			command("ls -al #2"),
		},
		servers: []server{
			{
				Host:  "a",
				Name:  "a",
				User:  "a",
				Fatal: "a",
			},
		},
		input: input{
			verifyFlag:  true,
			commandName: "a",
		},
	}

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		app.verifyAction()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		app.input.verifyFlag = false
		app.verifyAction()
		t.FailNow()
	}()

}
