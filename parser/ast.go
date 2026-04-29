package parser

import "fmt"

type TokenType int

const (
	TOK_SCENE TokenType = iota
	TOK_VAR
	TOK_LOCAL
	TOK_KEEP
	TOK_SET
	TOK_INPUT
	TOK_CHOICE
	TOK_BRANCH // Indented string inside a choice
	TOK_WHEN
	TOK_GO
	TOK_CALL
	TOK_FINISH
	TOK_END
	TOK_SAY
	TOK_TEXT
	TOK_DIALOGUE // A quoted line
	TOK_COMMENTS
)

func (t TokenType) String() string {
	switch t {
	case TOK_SCENE:
		return "SCENE"
	case TOK_VAR:
		return "VAR"
	case TOK_LOCAL:
		return "LOCAL"
	case TOK_KEEP:
		return "KEEP"
	case TOK_SET:
		return "SET"
	case TOK_INPUT:
		return "INPUT"
	case TOK_CHOICE:
		return "CHOICE"
	case TOK_BRANCH:
		return "BRANCH"
	case TOK_WHEN:
		return "WHEN"
	case TOK_GO:
		return "GO"
	case TOK_CALL:
		return "CALL"
	case TOK_FINISH:
		return "FINISH"
	case TOK_END:
		return "END"
	case TOK_SAY:
		return "SAY"
	case TOK_TEXT:
		return "TEXT"
	case TOK_DIALOGUE:
		return "DIALOGUE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

type Token struct {
	Type      TokenType
	Value     string
	Condition string
	Indent    int
	Line      int
}

type Story struct {
	Scenes []*SceneNode
}

type Node interface {
	NodeType() string
}

type SceneNode struct {
	Name string
	Body []Node
}

type VarNode struct {
	Name     string
	Value    any
	Lifetime TokenType
}

type SetNode struct {
	Raw string
}

type TextNode struct {
	Content    string
	IsDialogue bool
}

type ChoiceNode struct {
	Branches []*BranchNode
}

type BranchNode struct {
	Label     string
	Condition string
	Body      []Node
}

type JumpNode struct {
	Target string
	IsCall bool
}

type FinishNode struct{}
type EndNode struct{}

func (n *SceneNode) NodeType() string  { return "scene" }
func (n *VarNode) NodeType() string    { return "var" }
func (n *SetNode) NodeType() string    { return "set" }
func (n *TextNode) NodeType() string   { return "text" }
func (n *ChoiceNode) NodeType() string { return "choice" }
func (n *BranchNode) NodeType() string { return "branch" }
func (n *JumpNode) NodeType() string   { return "jump" }
func (n *FinishNode) NodeType() string { return "finish" }
func (n *EndNode) NodeType() string    { return "node" }
