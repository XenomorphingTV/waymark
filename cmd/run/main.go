package main

import (
	"fmt"
	"os"

	"github.com/XenomorphingTV/waymark/parser"
)

func main() {
	src, err := os.ReadFile("test_input.way")
	if err != nil {
		fmt.Println("error reading file:", err)
		os.Exit(1)
	}

	tokens, err := parser.Tokenize(string(src))
	if err != nil {
		fmt.Println("error tokenizing:", err)
		os.Exit(1)
	}

	for _, tok := range tokens {
		fmt.Printf("line %d indent %d type %v value %q condition %q\n",
			tok.Line, tok.Indent, tok.Type, tok.Value, tok.Condition)
	}

	story, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println("error parsing:", err)
		os.Exit(1)
	}

	for _, scene := range story.Scenes {
		fmt.Printf("scene: %s (%d nodes)\n", scene.Name, len(scene.Body))
	}
}
