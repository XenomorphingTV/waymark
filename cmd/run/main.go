package main

import (
	"fmt"
	"os"

	"github.com/XenomorphingTV/waymark/exporter"
	"github.com/XenomorphingTV/waymark/parser"
)

func main() {
	story, err := parser.ParseFile("test_input.way")
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	for _, scene := range story.Scenes {
		fmt.Printf("scene: %s (%d nodes)\n", scene.Name, len(scene.Body))
	}

	data, err := exporter.Export(story)
	if err != nil {
		fmt.Println("error exporting:", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}
