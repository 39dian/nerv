package cli

import (
	"testing"
	"regexp"
)


func TestAppInstall(t *testing.T) {

	//create
	cmd := &Cmd{
		Dir: "../../release/nerv/nerv-cli/bin",
		Cli:"./nerv-cli",
		Items:[]string{"topo", "create", "-t", "/demo/java/nerv-app-springboot.json", "-o", "nerv-app-springboot-1", "-n ../../../../test/cli/topo_test_inputs.json"},
	}

	var id string
	if out, err := cmd.Run(t); err != nil {
		t.Log(string(out))
		t.Errorf("%s", err.Error())
	} else {
		res := string(out)
		t.Log(res)
		regex := regexp.MustCompile(`.*([0-9]+),.*`)
		id = regex.FindStringSubmatch(res)[1]
	}

	//list
	cmd = &Cmd{
		Dir: "../../release/nerv/nerv-cli/bin",
		Cli:"./nerv-cli",
		Items:[]string{"topo", "list"},
	}

	if out, err := cmd.Run(t); err != nil {
		t.Log(string(out))
		t.Errorf("%s", err.Error())
	} else {
		t.Log(string(out))
	}

	//get
	cmd = &Cmd{
		Dir: "../../release/nerv/nerv-cli/bin",
		Cli:"./nerv-cli",
		Items:[]string{"topo", "get", "-i", id},
	}

	if out, err := cmd.Run(t); err != nil {
		t.Log(string(out))
		t.Errorf("%s", err.Error())
	} else {
		t.Log(string(out))
	}

	//install
	cmd = &Cmd{
		Dir: "../../release/nerv/nerv-cli/bin",
		Cli:"./nerv-cli",
		Items:[]string{"topo", "install", "-i", id},
	}

	if out, err := cmd.Run(t); err != nil {
		t.Log(string(out))
		t.Errorf("%s", err.Error())
	} else {
		t.Log(string(out))
	}

	//setup
	cmd = &Cmd{
		Dir: "../../release/nerv/nerv-cli/bin",
		Cli:"./nerv-cli",
		Items:[]string{"topo", "setup", "-i", id},
	}

	if out, err := cmd.Run(t); err != nil {
		t.Log(string(out))
		t.Errorf("%s", err.Error())
	} else {
		t.Log(string(out))
	}
}
