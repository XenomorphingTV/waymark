package main

import (
	"fmt"

	"github.com/XenomorphingTV/waymark/parser"
	"github.com/XenomorphingTV/waymark/runtime"
)

func main() {
	engine := runtime.New(&parser.Story{})

	// Seed some variables
	engine.SetVar("conversation_depth", 1)
	engine.SetVar("quest_active", true)
	engine.SetVar("player_strength", 10)
	engine.SetVar("gate_locked", false)

	tests := []struct {
		cond string
		want bool
	}{
		{"conversation_depth < 3", true},
		{"conversation_depth > 3", false},
		{"conversation_depth == 1", true},
		{"quest_active", true},
		{"not quest_active", false},
		{"not gate_locked", true},
		{"player_strength >= 10 and not gate_locked", true},
		{"player_strength >= 10 and gate_locked", false},
		{"conversation_depth < 3 or gate_locked", true},
		{"conversation_depth > 3 or quest_active", true},
		{"conversation_depth > 3 or gate_locked", false},
		{"(conversation_depth < 3 and quest_active) or gate_locked", true},
		{"conversation_depth < 3 and (quest_active or gate_locked)", true},
	}

	passed, failed := 0, 0
	for _, tt := range tests {
		got := engine.EvalCondition(tt.cond)
		status := "PASS"
		if got != tt.want {
			status = "FAIL"
			failed++
		} else {
			passed++
		}
		fmt.Printf("%s  %-50s got=%v want=%v\n", status, tt.cond, got, tt.want)
	}

	fmt.Printf("\n%d passed, %d failed\n", passed, failed)
}
