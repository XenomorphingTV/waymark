package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

// Parse builds a Story AST from a flat token stream. Top-level tokens must be SCENE definitions.
func Parse(tokens []Token) (*Story, error) {
	p := &Parser{tokens: tokens}
	story := &Story{}

	for !p.done() {
		if p.peek().Type != TOK_SCENE {
			return nil, fmt.Errorf("line %d: expected scene, got %v", p.peek().Line, p.peek().Type)
		}
		scene, err := p.parseScene()
		if err != nil {
			return nil, err
		}
		story.Scenes = append(story.Scenes, scene)
	}

	return story, nil
}

func (p *Parser) parseScene() (*SceneNode, error) {
	header := p.consume()
	scene := &SceneNode{Name: header.Value}
	baseIndent := header.Indent + 1

	for !p.done() && p.peek().Indent >= baseIndent {
		node, err := p.parseNode(baseIndent)
		if err != nil {
			return nil, err
		}
		scene.Body = append(scene.Body, node)
	}

	return scene, nil
}

func (p *Parser) parseNode(indent int) (Node, error) {
	tok := p.peek()

	// Skips tokens that are deeper than expected. Shouldn't happen if lexer is correct
	if tok.Indent != indent {
		return nil, fmt.Errorf("line %d: unexpected indent %d, expected %d", tok.Line, tok.Indent, indent)
	}

	switch tok.Type {
	case TOK_TEXT:
		p.consume()
		return &TextNode{Content: tok.Value, IsDialogue: false}, nil
	case TOK_DIALOGUE:
		p.consume()
		return &TextNode{Content: tok.Value, IsDialogue: true}, nil
	case TOK_VAR, TOK_LOCAL, TOK_KEEP:
		return p.parseValDecl()
	case TOK_SET:
		p.consume()
		return &SetNode{Raw: tok.Value}, nil
	case TOK_CHOICE:
		return p.parseChoice(indent)
	case TOK_GO:
		p.consume()
		return &JumpNode{Target: tok.Value, IsCall: false}, nil
	case TOK_CALL:
		p.consume()
		return &JumpNode{Target: tok.Value, IsCall: true}, nil
	case TOK_FINISH:
		p.consume()
		return &FinishNode{}, nil
	case TOK_END:
		p.consume()
		return &EndNode{}, nil
	}

	return nil, fmt.Errorf("line %d: unexpected token %v", tok.Line, tok.Type)
}

func (p *Parser) parseValDecl() (*VarNode, error) {
	tok := p.consume()

	parts := strings.SplitN(tok.Value, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("line %d: malformed variable declaration %q", tok.Line, tok.Value)
	}

	return &VarNode{
		Name:     strings.TrimSpace(parts[0]),
		Value:    strings.TrimSpace(parts[1]),
		Lifetime: tok.Type,
	}, nil
}

func (p *Parser) parseChoice(indent int) (*ChoiceNode, error) {
	p.consume()

	choice := &ChoiceNode{}
	branchIndent := indent + 1

	for !p.done() && p.peek().Indent == branchIndent {
		tok := p.peek()

		// Both BRANCH and DIALOGUE at this indent are choice options
		// DIALOGUE here means a quoted string with no `when`
		if tok.Type != TOK_BRANCH && tok.Type != TOK_DIALOGUE {
			break
		}

		branch, err := p.parseBranch(branchIndent)
		if err != nil {
			return nil, err
		}
		choice.Branches = append(choice.Branches, branch)
	}

	if len(choice.Branches) == 0 {
		return nil, fmt.Errorf("line %d: choice block has no branches", p.peek().Line)
	}

	return choice, nil
}

func (p *Parser) parseBranch(indent int) (*BranchNode, error) {
	tok := p.consume() // BRANCH or DIALOGUE

	branch := &BranchNode{
		Label:     tok.Value,
		Condition: tok.Condition,
	}

	bodyIndent := indent + 1

	for !p.done() && p.peek().Indent >= bodyIndent {
		node, err := p.parseNode(bodyIndent)
		if err != nil {
			return nil, err
		}

		branch.Body = append(branch.Body, node)
	}

	return branch, nil
}

func (p *Parser) peek() Token {
	return p.tokens[p.pos]
}

func (p *Parser) consume() Token {
	t := p.tokens[p.pos]
	p.pos++
	return t
}

func (p *Parser) done() bool {
	return p.pos >= len(p.tokens)
}
