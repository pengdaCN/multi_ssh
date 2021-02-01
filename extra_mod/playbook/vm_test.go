package playbook

import "testing"

func TestTmpl(t *testing.T) {
	err := VM.DoFile("./test_template.lua")
	if err != nil {
		t.Fatal(err)
	}
}
