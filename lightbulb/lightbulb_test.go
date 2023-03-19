package lightbulb

import "testing"

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

func TestFindAllblocks(t *testing.T) {
	text, err := LoadFromFile("good.md")
	if err != nil {
		t.Error(err)
	}
	blocks, err := findAllBlocks(text)
	if err != nil {
		t.Error(err)
		return
	}
	if len(blocks) != 2 {
		t.Error("Expected 2 blocks, got ", len(blocks))
		return
	}

	sections := []string{
		`<!-- lightbulb:createFile name:dateFile path:./date.sh tags:one,two,three mode:0700 -->
` + "```" + `shell
#!/bin/bash

echo "The current date in UTC is $(date -u)."
` + "```",
		`<!-- lightbulb:runShell name:runDate shell:bash -->
` + "```" + `console
./date.sh
` + "```",
	}

	for i, section := range sections {
		if blocks[i] != section {
			t.Error("Block ", i, " doesn't match expected")
		}
	}
}

func TestGetBlockAction(t *testing.T) {
	block := Block{}
	text := `<!-- lightbulb:createFile name:dateFile path:date.sh tags:one,two,three mode:0700 -->`
	err := getBlockAction(text, &block)
	if err != nil {
		t.Error(err)
	}
	if block.Action != "createFile" {
		t.Error("Expected action createFile, got ", block.Action)
	}
	text = `<!-- lightbulb:runShell name:runDate shell:bash -->`
	block = Block{}
	err = getBlockAction(text, &block)
	if err != nil {
		t.Error(err)
	}
	if block.Action != "runShell" {
		t.Error("Expected action runShell, got ", block.Action)
	}
	text = `<!-- lightbulb: name:runDate shell:bash -->`
	block = Block{}
	err = getBlockAction(text, &block)
	if err == nil {
		t.Error("Expected error, got ", block.Action)
	}
	if block.Action != "" {
		t.Error("Expected action '', got ", block.Action)
	}

}

func TestGetBlockParameters(t *testing.T) {
	block := Block{}
	text := `<!-- lightbulb:createFile name:dateFile path:date.sh tags:one,two,three mode:0700 -->`
	err := getBlockParameters(text, &block)
	if err != nil {
		t.Error(err)
	}
	if block.Name != "dateFile" {
		t.Error("Expected name dateFile, got ", block.Name)
	}
	if block.Path != "date.sh" {
		t.Error("Expected path date.sh, got ", block.Path)
	}
	if !arrayMatch(block.Tags, []string{"one", "two", "three"}) {
		t.Error("Expected tags [one two three], got ", block.Tags)
	}
	if block.Mode != "0700" {
		t.Error("Expected mode 0700, got ", block.Mode)
	}

}

func TestGetBlockCode(t *testing.T) {
	block := Block{}
	text := "```" + `shell
#!/bin/bash

echo "The current date in UTC is $(date -u)."
` + "```"
	err := getBlockCode(text, &block)
	if err != nil {
		t.Error(err)
	}
	if block.Code != `#!/bin/bash

echo "The current date in UTC is $(date -u)."
` {
		t.Errorf("Expected code, got '%s'", block.Code)
	}

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
	if blocks[0].Action != "createFile" {
		t.Error("Expected block.Action createFile, got ", blocks[0].Action)
	}
	if blocks[0].Name != "dateFile" {
		t.Error("Expected block.Name dateFile, got ", blocks[0].Name)
	}
	if blocks[0].Path != "./date.sh" {
		t.Error("Expected block.Path ./date.sh, got ", blocks[0].Path)
	}
	if !arrayMatch(blocks[0].Tags, []string{"one", "two", "three"}) {
		t.Error("Expected block.Tags [one two three], got ", blocks[0].Tags)
	}
	if blocks[0].Mode != "0700" {
		t.Error("Expected block.Mode 0700, got ", blocks[0].Mode)
	}
	if blocks[0].Code != `#!/bin/bash

echo "The current date in UTC is $(date -u)."
` {
		t.Errorf("Expected block.Code, got '%s'", blocks[0].Code)
	}

	if blocks[1].Action != "runShell" {
		t.Error("Expected block.Action runShell, got ", blocks[1].Action)
	}
	if blocks[1].Name != "runDate" {
		t.Error("Expected block.Name runDate, got ", blocks[1].Name)
	}
	if blocks[1].Shell != "bash" {
		t.Error("Expected block.Shell bash, got ", blocks[1].Shell)
	}

}

func TestGoodRun(t *testing.T) {
	text, err := LoadFromFile("good.md")
	if err != nil {
		t.Error(err)
		return
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
	err = Run(blocks)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBadRun(t *testing.T) {
	text, err := LoadFromFile("error.md")
	if err != nil {
		t.Error(err)
		return
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
	err = Run(blocks)
	if err == nil {
		t.Error("Expected error, got nil")
		return
	}
}
