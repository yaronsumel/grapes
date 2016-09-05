package main

import (
	"testing"
)

var demoInputOk = input{
	asyncFlag:   asyncFlag(false),
	configPath:  configPath("/demo"),
	keyPath:     keyPath("/demo"),
	serverGroup: serverGroup("demoServerGroup"),
	commandName: commandName("demo"),
	verifyFlag:  verifyFlag(false),
}

var demoInputNotOk = input{
	asyncFlag:   asyncFlag(false),
	configPath:  configPath(""),
	keyPath:     keyPath(""),
	serverGroup: serverGroup(""),
	verifyFlag:  verifyFlag(true),
}

func TestParse(t *testing.T) {
	in := input{}
	in.parse()
}

func TestValidateInput(t *testing.T) {
	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.validate()
	}()
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputNotOk.validate()
		t.FailNow()
	}()

}

func TestValidateConfigPath(t *testing.T) {
	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.configPath.validate()
	}()
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputNotOk.configPath.validate()
		t.FailNow()
	}()
}

func TestValidateKeyPath(t *testing.T) {
	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.keyPath.validate()
	}()
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputNotOk.keyPath.validate()
		t.FailNow()
	}()
}

func TestValidateServerGroup(t *testing.T) {
	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.serverGroup.validate()
	}()
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputNotOk.serverGroup.validate()
		t.FailNow()
	}()
}

func TestValidateCommandName(t *testing.T) {
	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.commandName.validate()
	}()
	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputNotOk.commandName.validate()
		t.FailNow()
	}()
}

func TestVerifyAction(t *testing.T) {

	s := servers{
		server{
			Name: "testN",
			Host: "testH",
			User: "testU",
		},
	}

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		demoInputOk.verifyAction(s)
		t.FailNow()
	}()

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		demoInputOk.verifyFlag = true
		demoInputOk.verifyAction(s)
	}()

}
