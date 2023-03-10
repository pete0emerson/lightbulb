package lightbulb

import (
	"testing"
)

func arrayMatch(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i, _ := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	text, err := LoadFromFile("good.md")
	if err != nil {
		t.Error(err)
	}
	blocks, err := Parse(text)
	if err != nil {
		t.Error(err)
		return
	}
	if len(blocks) != 2 {
		t.Error("Expected 2 blocks, got ", len(blocks))
		return
	}

	goodCompare := []Block{
		{
			Action:      "createFile",
			Name:        "dateFile",
			Tags:        []string{"one", "two", "three"},
			Path:        "date.sh",
			Mode:        "0700",
			Command:     "",
			Shell:       "",
			Set:         "",
			ExitOnError: false,
			Keys:        []string{},
			Prompt:      false,
			Secret:      false,
			Persist:     false,
			Sensitive:   false,
		},
		{
			Action:      "runShell",
			Name:        "runDate",
			Tags:        []string{},
			Path:        "",
			Mode:        "",
			Command:     "",
			Shell:       "bash",
			Set:         "",
			ExitOnError: false,
			Keys:        []string{},
			Prompt:      false,
			Secret:      false,
			Persist:     false,
			Sensitive:   false,
		},
	}

	for i, block := range goodCompare {
		if block.Action != blocks[i].Action {
			t.Error("Expected block.Action ", block.Action, " got ", blocks[i].Action)
		}
		if block.Name != blocks[i].Name {
			t.Error("Expected block.Name ", block.Name, " got ", blocks[i].Name)
		}
		if !arrayMatch(block.Tags, blocks[i].Tags) {
			t.Error("Expected block.Tags ", block.Tags, " got ", blocks[i].Tags)
		}
		if block.Path != blocks[i].Path {
			t.Error("Expected block.Path ", block.Path, " got ", blocks[i].Path)
		}
		if block.Mode != blocks[i].Mode {
			t.Error("Expected block.Mode ", block.Mode, " got ", blocks[i].Mode)
		}
		if block.Command != blocks[i].Command {
			t.Error("Expected block.Command ", block.Command, " got ", blocks[i].Command)
		}
		if block.Shell != blocks[i].Shell {
			t.Error("Expected block.Shell ", block.Shell, " got ", blocks[i].Shell)
		}
		if block.Set != blocks[i].Set {
			t.Error("Expected block.Set ", block.Set, " got ", blocks[i].Set)
		}
		if block.Prompt != blocks[i].Prompt {
			t.Error("Expected block.Prompt ", block.Prompt, " got ", blocks[i].Prompt)
		}
		if !arrayMatch(block.Keys, blocks[i].Keys) {
			t.Error("Expected block.Keys ", block.Keys, " got ", blocks[i].Keys)
		}
		if block.Prompt != blocks[i].Prompt {
			t.Error("Expected block.Prompt ", block.Prompt, " got ", blocks[i].Prompt)
		}
		if block.Secret != blocks[i].Secret {
			t.Error("Expected block.Secret ", block.Secret, " got ", blocks[i].Secret)
		}
		if block.Persist != blocks[i].Persist {
			t.Error("Expected block.Persist ", block.Persist, " got ", blocks[i].Persist)
		}
		if block.Sensitive != blocks[i].Sensitive {
			t.Error("Expected block.Sensitive ", block.Sensitive, " got ", blocks[i].Sensitive)
		}

	}

}
