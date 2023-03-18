package lightbulb

import (
	"errors"
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
	Code        string
	Shell       string
	Set         string
	ExitOnError bool
	Keys        []string
	Prompt      string
	Secret      bool
	Persist     bool
	Sensitive   bool
}

// findAllBlocks finds all raw text blocks in a markdown string
func findAllBlocks(content string) ([]string, error) {
	re := regexp.MustCompile("(?s)<!-- lightbulb:.*? -->.?```.*?```")
	log.Debugf("Finding all blocks matching %s", re.String())
	rawBlocks := re.FindAllString(content, -1)
	return rawBlocks, nil
}

// getBlockAction returns the Lightbulb action from an HTML comment
func getBlockAction(text string, block *Block) error {
	re := regexp.MustCompile(`<!-- lightbulb:(\S+).*?-->`)
	log.Debugf("Finding block action matching %s", re.String())
	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return errors.New("No matches in block found")
	} else {
		block.Action = matches[1]
	}
	return nil
}

func setBlockDefaults(block *Block) {
	block.Tags = []string{"all"}
	block.Mode = "0700"
	block.Shell = "bash"
	block.ExitOnError = true
	block.Prompt = "missing"
	block.Secret = false
	block.Persist = true
	block.Sensitive = false
}

// getBlockParameters returns the Lightbulb parameters from an HTML comment
func getBlockParameters(text string, block *Block) error {
	setBlockDefaults(block)
	re := regexp.MustCompile(`<!-- lightbulb:\S+ (.*)-->`)
	log.Debugf("Finding block parameters matching %s", re.String())
	parameterMatches := re.FindStringSubmatch(text)
	if len(parameterMatches) == 0 {
		log.Infof("No parameters found in block")
		return nil
	}
	parameters := parameterMatches[1]
	parameters = strings.TrimSpace(parameters)
	log.Debugf("Extracted parameters: %s", parameters)
	re = regexp.MustCompile(`(\S+):(\S+)`)
	log.Debugf("Finding parameters matching %s", re.String())
	matches := re.FindAllStringSubmatch(parameters, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		switch key {
		case "name":
			block.Name = value
		case "tags":
			block.Tags = strings.Split(value, ",")
		case "path":
			block.Path = value
		case "mode":
			block.Mode = value
		case "shell":
			block.Shell = value
		case "set":
			block.Set = value
		case "exit_on_error":
			if value == "true" {
				block.ExitOnError = true
			}
		case "keys":
			block.Keys = strings.Split(value, ",")
		case "prompt":
			block.Prompt = value
		case "secret":
			if value == "true" {
				block.Secret = true
			}
		case "persist":
			if value == "true" {
				block.Persist = true
			}
		case "sensitive":
			if value == "true" {
				block.Sensitive = true
			}
		default:
			return fmt.Errorf("unknown parameter: %s", key)
		}

	}
	return nil
}

// getBlockCode returns the code block from a markdown string
func getBlockCode(text string, block *Block) error {
	re := regexp.MustCompile("(?s)```[^\n]*\n(.*?)```")
	log.Debugf("Finding block code matching %s", re.String())
	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return errors.New("no matches in block found")
	} else {
		block.Code = matches[1]
	}
	return nil

}

// processTextBlock processes a raw text block into a Block
func processTextBlock(text string) (Block, error) {
	log.Debugf("Processing text block:\n%s", text)
	block := Block{}
	err := getBlockAction(text, &block)
	if err != nil {
		return block, err
	}

	err = getBlockParameters(text, &block)
	if err != nil {
		return block, err
	}

	err = getBlockCode(text, &block)
	if err != nil {
		return block, err
	}

	return block, nil
}

// Parse parses a markdown string into a slice of Blocks
func Parse(content string) ([]Block, error) {
	log.Info("Parsing markdown")
	log.Debugf("Markdown:\n%s", content)

	rawBlocks, err := findAllBlocks(content)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d blocks", len(rawBlocks))

	blocks := []Block{}
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
	log.Infof("Loading markdown from file: %s", fileName)

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	log.Debugf("File contents:\n%s", content)
	return string(content), nil
}

// LoadFromURL loads a markdown file from a URL into a string
func LoadFromURL(url string) (string, error) {
	log.Debugf("Loading markdown from URL: %s", url)
	return "", nil
}
