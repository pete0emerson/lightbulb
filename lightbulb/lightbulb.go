package lightbulb

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Block struct {
	Action      string
	Name        string
	Tags        []string
	Path        string
	Mode        string
	Command     string
	Shell       string
	Set         string
	ExitOnError bool
	Keys        []string
	Prompt      bool
	Secret      bool
	Persist     bool
	Sensitive   bool
}

// findAllBlocks finds all raw text blocks in a markdown string
func findAllBlocks(content string) ([]string, error) {
	log.Debug("Finding all blocks")
	re := regexp.MustCompile("(?s)<!-- lightbulb:.*? -->.?```.*?```")
	rawBlocks := re.FindAllString(content, -1)

	return rawBlocks, nil

}

// processTextBlock processes a raw text block into a Block
func processTextBlock(text string) (Block, error) {
	log.Debug("Processing text block")
	block := Block{}

	re := regexp.MustCompile("<!-- lightbulb:(" + `\S+` + ").*? -->")
	match := re.FindStringSubmatch(text)

	lightbulbCommand := match[0]
	block.Action = match[1]
	log.Debugf("Block Action: %s", block.Action)

	for _, param := range []string{"name", "tags", "path", "mode", "command", "shell", "set", "exitOnError", "keys", "prompt", "secret", "persist", "sensitive"} {
		re = regexp.MustCompile(param + ":(" + `\S+` + ")")
		match = re.FindStringSubmatch(lightbulbCommand)

		// Required parameters
		switch block.Action {
		case "createFile":
			if len(match) == 0 && (param == "name" || param == "path") {
				return block, fmt.Errorf("no required parameter '%s' found in lightbulb '%s' action", param, block.Action)
			}
		case "runShell":
			if len(match) == 0 && param == "name" {
				return block, fmt.Errorf("no required parameter '%s' found in lightbulb '%s' action", param, block.Action)
			}
		case "setEnvironmentVars":
			if len(match) == 0 && (param == "name" || param == "keys") {
				return block, fmt.Errorf("no required parameter '%s' found in lightbulb '%s' action", param, block.Action)
			}
		}

		if len(match) > 0 {
			switch param {
			case "name":
				block.Name = match[1]
			case "tags":
				block.Tags = strings.Split(match[1], ",")
			case "path":
				block.Path = match[1]
			case "mode":
				block.Mode = match[1]
			case "command":
				block.Command = match[1]
			case "shell":
				block.Shell = match[1]
			case "set":
				block.Set = match[1]
			case "exitOnError":
				if match[1] == "true" {
					block.ExitOnError = true
				} else if match[1] == "false" {
					block.ExitOnError = false
				} else {
					return block, fmt.Errorf("parameter '%s' found in lightbulb '%s' action is not true or false", param, block.Action)
				}
			case "keys":
				block.Keys = strings.Split(match[1], ",")
			case "prompt":
				if match[1] == "true" {
					block.Prompt = true
				} else if match[1] == "false" {
					block.Prompt = false
				} else {
					return block, fmt.Errorf("parameter '%s' found in lightbulb '%s' action is not true or false", param, block.Action)
				}
			case "secret":
				if match[1] == "true" {
					block.Secret = true
				} else if match[1] == "false" {
					block.Secret = false
				} else {
					return block, fmt.Errorf("parameter '%s' found in lightbulb '%s' action is not true or false", param, block.Action)
				}
			case "persist":
				if match[1] == "true" {
					block.Persist = true
				} else if match[1] == "false" {
					block.Persist = false
				} else {
					return block, fmt.Errorf("parameter '%s' found in lightbulb '%s' action is not true or false", param, block.Action)
				}
			case "sensitive":
				if match[1] == "true" {
					block.Sensitive = true
				} else if match[1] == "false" {
					block.Sensitive = false
				} else {
					return block, fmt.Errorf("parameter '%s' found in lightbulb '%s' action is not true or false", param, block.Action)
				}

			}
		}

	}

	return block, nil
}

// Parse parses a markdown string into a slice of Blocks
func Parse(content string) ([]Block, error) {
	rawBlocks, err := findAllBlocks(content)
	if err != nil {
		return nil, err
	}

	var blocks []Block

	for _, rawBlock := range rawBlocks {
		block, err := processTextBlock(rawBlock)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil

}

// LoadFromFile loads a markdown file into a string
func LoadFromFile(fileName string) (string, error) {
	log.Debugf("Loading markdown from file: %s", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// LoadFromURL loads a markdown file from a URL into a string
func LoadFromURL(url string) (string, error) {
	log.Debugf("Loading markdown from URL: %s", url)
	return "", nil
}
