package parser

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{
			name:  "simple scene",
			input: "scene tavern",
			want: []Token{
				{Type: TOK_SCENE, Value: "tavern", Indent: 0, Line: 1},
			},
		},
		{
			name:  "scene with text body",
			input: "scene tavern\n\tThe barkeep eyes you.",
			want: []Token{
				{Type: TOK_SCENE, Value: "tavern", Indent: 0, Line: 1},
				{Type: TOK_TEXT, Value: "The barkeep eyes you.", Indent: 1, Line: 2},
			},
		},
		{
			name:  "skips blank lines",
			input: "scene tavern\n\n\tSome text.",
			want: []Token{
				{Type: TOK_SCENE, Value: "tavern", Indent: 0, Line: 1},
				{Type: TOK_TEXT, Value: "Some text.", Indent: 1, Line: 3},
			},
		},
		{
			name:  "skips comments",
			input: "scene tavern\n\t# this is a comment\n\tSome text.",
			want: []Token{
				{Type: TOK_SCENE, Value: "tavern", Indent: 0, Line: 1},
				{Type: TOK_TEXT, Value: "Some text.", Indent: 1, Line: 3},
			},
		},
		{
			name:  "windows line endings",
			input: "scene tavern\r\n\tSome text.\r\n",
			want: []Token{
				{Type: TOK_SCENE, Value: "tavern", Indent: 0, Line: 1},
				{Type: TOK_TEXT, Value: "Some text.", Indent: 1, Line: 2},
			},
		},
		{
			name:    "space indent returns error",
			input:   "scene tavern\n    Some text.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d tokens, want %d\ngot:  %+v\nwant: %+v",
					len(got), len(tt.want), got, tt.want)
			}

			for i, tok := range got {
				if tok != tt.want[i] {
					t.Errorf("token %d:\n  got  %+v\n  want %+v", i, tok, tt.want[i])
				}
			}
		})
	}
}

// TODO: Finish implementing the lexer tests
