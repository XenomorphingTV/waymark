package runtime

import "github.com/XenomorphingTV/waymark/parser"

type Engine struct {
	story     *parser.Story
	vars      map[string]any
	callStack []frame
	pos       cursor
}

type frame struct {
	scene    *parser.SceneNode
	pos      int
	locals   map[string]any
	keepVars map[string]bool
}

type cursor struct {
	scene *parser.SceneNode
	pos   int
}

type State struct {
	Lines   []Line
	Choices []Choice
	Done    bool
}

type Line struct {
	Content    string
	IsDialogue bool
}

type Choice struct {
	Index int
	ID    string
	Label string
}

type conditionParser struct {
	input string
	pos   int
}
