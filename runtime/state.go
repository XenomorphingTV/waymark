package runtime

import "github.com/XenomorphingTV/waymark/parser"

type Engine struct {
	story     *parser.Story
	vars      map[string]any
	callbacks []callback
	pos       cursor
}

type callback struct {
	scene     *parser.SceneNode
	pos       int
	localvars map[string]any
}

type cursor struct {
	scene *parser.SceneNode
	pos   int
}

type State struct {
	Lines   []Line
	Choices []Choice
}

type Line struct {
	Content    string
	IsDialogue bool
}

type Choice struct {
	Index int
	Label string
}
