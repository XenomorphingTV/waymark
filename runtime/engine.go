package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/XenomorphingTV/waymark/parser"
)

// TEMP TEST FUNCTIONS -------

func (e *Engine) SetVar(name string, val any) {
	e.vars[name] = val
}

func (e *Engine) EvalCondition(cond string) bool {
	return e.evalCondition(cond)
}

// ---------------------------

func New(story *parser.Story) *Engine {
	return &Engine{
		story: story,
		vars:  make(map[string]any),
	}
}

func (e *Engine) Start(sceneName string) (*State, error) {
	scene, ok := e.story.Scene(sceneName)
	if !ok {
		return nil, fmt.Errorf("scene %q not found", sceneName)
	}

	e.callStack = []frame{{
		scene:    scene,
		pos:      0,
		locals:   make(map[string]any),
		keepVars: make(map[string]bool),
	}}

	return e.run()
}

func (e *Engine) Choose(index int) (*State, error) {
	if len(e.callStack) == 0 {
		return nil, fmt.Errorf("engine not started")
	}

	// Slice of frames. Last one is always the one currently executing
	currentFrame := &e.callStack[len(e.callStack)-1]

	// Get the ChoiceNode we are paused at. We type assert it as one
	node := currentFrame.scene.Body[currentFrame.pos].(*parser.ChoiceNode)

	// Collect available branches (only those which have passing "when" conditions)
	available := e.availableBranches(node)

	if index < 0 || index >= len(available) {
		return nil, fmt.Errorf("invalid choice at index %d", index)
	}

	branch := available[index]

	e.callStack[len(e.callStack)-1].pos++

	e.callStack = append(e.callStack, frame{
		scene:    &parser.SceneNode{Body: branch.Body},
		pos:      0,
		locals:   make(map[string]any),
		keepVars: make(map[string]bool),
	})

	return e.run()
}

func (e *Engine) run() (*State, error) {
	return nil, nil
}

func (e *Engine) setVar(name string, val any) {
	// Check locals top to bottom
	for i := len(e.callStack) - 1; i >= 0; i-- {
		if _, ok := e.callStack[i].locals[name]; ok {
			e.callStack[i].locals[name] = val
			return
		}
	}
	// Fall back to globals
	e.vars[name] = val
}

func (e *Engine) getVar(name string) any {
	if val, ok := e.localVar(name); ok {
		return val
	}
	return e.vars[name]
}

func (e *Engine) availableBranches(node *parser.ChoiceNode) []*parser.BranchNode {
	var available []*parser.BranchNode
	for _, branch := range node.Branches {
		if branch.Condition == "" || e.evalCondition(branch.Condition) {
			available = append(available, branch)
		}
	}
	return available
}

func (e *Engine) evalCondition(cond string) bool {
	p := &conditionParser{input: strings.TrimSpace(cond)}
	return e.evalOr(p)
}

func (e *Engine) evalOr(p *conditionParser) bool {
	left := e.evalAnd(p)
	for p.peek("or") {
		p.consume("or")
		right := e.evalAnd(p)
		left = left || right
	}
	return left
}

func (e *Engine) evalAnd(p *conditionParser) bool {
	left := e.evalNot(p)
	for p.peek("and") {
		p.consume("and")
		right := e.evalNot(p)
		left = left && right
	}
	return left
}

func (e *Engine) evalNot(p *conditionParser) bool {
	if p.peek("not") {
		p.consume("not")
		return !e.evalNot(p)
	}
	return e.evalComparison(p)
}

func (e *Engine) evalComparison(p *conditionParser) bool {
	left := e.evalPrimary(p)
	p.skipSpace()

	for _, op := range []string{">=", "<=", "!=", "==", ">", "<"} {
		if strings.HasPrefix(p.rest(), op) {
			p.pos += len(op)
			p.skipSpace()
			right := e.evalPrimary(p)
			result := compare(left, op, right)
			fmt.Printf("DEBUG: %v %s %v = %v\n", left, op, right, result)
			return result
		}
	}

	switch v := left.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != ""
	default:
		return left != nil
	}
}

func (e *Engine) evalPrimary(p *conditionParser) any {
	p.skipSpace()

	// Parenthesised expression
	if p.pos < len(p.input) && p.input[p.pos] == '(' {
		p.pos++ // consume "("
		val := e.evalOr(p)
		p.skipSpace()
		if p.pos < len(p.input) && p.input[p.pos] == ')' {
			p.pos++ // consume ")"
		}

		return val
	}

	return e.resolveValue(p.readToken())
}

// Helpers

func (p *conditionParser) rest() string {
	return p.input[p.pos:]
}

func (p *conditionParser) skipSpace() {
	for p.pos < len(p.input) && p.input[p.pos] == ' ' {
		p.pos++
	}
}

func (p *conditionParser) peek(keyword string) bool {
	p.skipSpace()
	rest := p.rest()
	if !strings.HasPrefix(rest, keyword) {
		return false
	}

	// Make sure it is a whole word, not a prefix of something else
	after := len(keyword)
	if after < len(rest) && rest[after] != ' ' && rest[after] != ')' {
		return false
	}
	return true
}

func (p *conditionParser) consume(keyword string) {
	p.skipSpace()
	p.pos += len(keyword)
}

func (p *conditionParser) readToken() string {
	p.skipSpace()
	start := p.pos
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if c == ' ' || c == ')' || c == '(' {
			break
		}

		// Stop before operators
		rest := p.rest()
		if strings.HasPrefix(rest, "==") || strings.HasPrefix(rest, "!=") ||
			strings.HasPrefix(rest, ">=") || strings.HasPrefix(rest, "<=") ||
			strings.HasPrefix(rest, ">") || strings.HasPrefix(rest, "<") {
			break
		}
		p.pos++
	}
	return p.input[start:p.pos]
}

func (e *Engine) resolveValue(s string) any {
	// String literal
	if strings.HasPrefix(s, `"`) {
		return strings.Trim(s, `"`)
	}
	// Bool literal
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	// Int literal
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	// Float literal
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// Variable lookup — check locals first, then globals
	if val, ok := e.localVar(s); ok {
		return val
	}
	if val, ok := e.vars[s]; ok {
		return val
	}
	// Unknown variable — return nil
	return nil
}

func (e *Engine) localVar(name string) (any, bool) {
	for i := len(e.callStack) - 1; i >= 0; i-- {
		if val, ok := e.callStack[i].locals[name]; ok {
			return val, true
		}
	}
	return nil, false
}

func compare(left any, op string, right any) bool {
	// numeric comparison
	l, lok := toFloat(left)
	r, rok := toFloat(right)
	if lok && rok {
		switch op {
		case "==":
			return l == r
		case "!=":
			return l != r
		case ">":
			return l > r
		case "<":
			return l < r
		case ">=":
			return l >= r
		case "<=":
			return l <= r
		}
	}

	// String/bool equality
	switch op {
	case "==":
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
	case "!=":
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right)
	}

	return false
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case float64:
		return n, true
	}
	return 0, false
}
