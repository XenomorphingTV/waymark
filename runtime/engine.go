package runtime

import "github.com/XenomorphingTV/waymark/parser"

func New(story *parser.Story) *Engine {
	return &Engine{
		story: story,
		vars:  make(map[string]any),
	}
}

func (e *Engine) Start(scene string) (*State, error) {
}
