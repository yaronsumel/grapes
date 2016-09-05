package main

import (
	"testing"
)

//func TestMain(t *testing.M){
//
//}

func TestRunApp(t *testing.T) {

	app := grape{}

	app.input.asyncFlag = true
	app.input.verifyFlag = true
	app.input.serverGroup = "prod"
	app.input.commandName = "date"

	app.config.set("testFiles/demo_config.yml")
	app.ssh.setKey(keyPath("testFiles/id_rsa"))

	//should not panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		app.runApp()
	}()

	//should panic
	func() {
		defer func() {
			if err := recover(); err != nil {
				recover()
			}
		}()
		app.input.asyncFlag = false
		app.config.Servers["prod"][0].Host = "notworking.local"
		app.runApp()
		t.FailNow()
	}()

}

func TestPrint(t *testing.T) {

	s := server{
		Host: "sdf.org:22",
		User: "new",
		Name: "public",
	}

	x := grapeCommandStd{
		Command: "a",
		Std: std{
			Err: "",
			Out: "",
		},
	}

	so := serverOutput{
		Server: s,
		Output: []*grapeCommandStd{
			&x,
		},
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				t.FailNow()
			}
		}()
		so.print()
	}()
}
