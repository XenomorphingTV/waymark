package exporter

import (
	"encoding/json"
	"fmt"

	"github.com/XenomorphingTV/waymark/parser"
)

type StoryJSON struct {
	Scenes []SceneJSON `json:"scenes"`
}

type SceneJSON struct {
	Name string `json:"name"`
	Body []any  `json:"body"`
}

type VarNodeJSON struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    any    `json:"value"`
	Lifetime string `json:"lifetime"`
}

type SetNodeJSON struct {
	Type string `json:"type"`
	Raw  string `json:"raw"`
}

type TextNodeJSON struct {
	Type       string `json:"type"`
	Content    string `json:"content"`
	IsDialogue bool   `json:"is_dialogue"`
}

type ChoiceNodeJSON struct {
	Type     string            `json:"type"`
	Branches []*BranchNodeJSON `json:"branches"`
}

type BranchNodeJSON struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Label     string `json:"label"`
	Condition string `json:"condition"`
	Body      []any  `json:"body"`
}

type JumpNodeJSON struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	IsCall bool   `json:"is_call"`
}

type FinishNodeJSON struct {
	Type string `json:"type"`
}

type EndNodeJSON struct {
	Type string `json:"type"`
}

func Export(story *parser.Story) ([]byte, error) {
	var scenes []SceneJSON

	for _, scene := range story.Scenes {
		exported, err := exportScene(scene)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, exported)
	}

	return json.MarshalIndent(StoryJSON{Scenes: scenes}, "", "  ")
}

func exportScene(scene *parser.SceneNode) (SceneJSON, error) {
	var body []any

	for _, node := range scene.Body {
		exported, err := exportNode(node)
		if err != nil {
			return SceneJSON{}, err
		}
		body = append(body, exported)
	}

	return SceneJSON{Name: scene.Name, Body: body}, nil
}

func exportNode(node parser.Node) (any, error) {
	switch v := node.(type) {
	case *parser.VarNode:
		return exportVar(v), nil
	case *parser.SetNode:
		return exportSet(v), nil
	case *parser.TextNode:
		return exportText(v), nil
	case *parser.ChoiceNode:
		return exportChoice(v)
	case *parser.JumpNode:
		return exportJump(v), nil
	case *parser.FinishNode:
		return exportFinish(v), nil
	case *parser.EndNode:
		return exportEnd(v), nil
	default:
		return nil, fmt.Errorf("unknown node type %T", node)
	}
}

func exportVar(node *parser.VarNode) VarNodeJSON {
	return VarNodeJSON{
		Type:     node.NodeType(),
		Name:     node.Name,
		Value:    node.Value,
		Lifetime: node.Lifetime.String(),
	}
}

func exportSet(node *parser.SetNode) SetNodeJSON {
	return SetNodeJSON{
		Type: node.NodeType(),
		Raw:  node.Raw,
	}
}

func exportText(node *parser.TextNode) TextNodeJSON {
	return TextNodeJSON{
		Type:       node.NodeType(),
		Content:    node.Content,
		IsDialogue: node.IsDialogue,
	}
}

func exportChoice(node *parser.ChoiceNode) (ChoiceNodeJSON, error) {
	var branches []*BranchNodeJSON

	for _, branch := range node.Branches {
		exported, err := exportBranch(branch)
		if err != nil {
			return ChoiceNodeJSON{}, err
		}
		branches = append(branches, &exported)
	}

	return ChoiceNodeJSON{
		Type:     node.NodeType(),
		Branches: branches,
	}, nil
}

func exportBranch(node *parser.BranchNode) (BranchNodeJSON, error) {
	var body []any

	for _, inner_node := range node.Body {
		exported, err := exportNode(inner_node)
		if err != nil {
			return BranchNodeJSON{}, err
		}
		body = append(body, exported)
	}

	return BranchNodeJSON{
		Type:      node.NodeType(),
		ID:        node.ID,
		Label:     node.Label,
		Condition: node.Condition,
		Body:      body,
	}, nil
}

func exportJump(node *parser.JumpNode) JumpNodeJSON {
	return JumpNodeJSON{
		Type:   node.NodeType(),
		Target: node.Target,
		IsCall: node.IsCall,
	}
}

func exportFinish(node *parser.FinishNode) FinishNodeJSON {
	return FinishNodeJSON{
		Type: node.NodeType(),
	}
}

func exportEnd(node *parser.EndNode) EndNodeJSON {
	return EndNodeJSON{
		Type: node.NodeType(),
	}
}
