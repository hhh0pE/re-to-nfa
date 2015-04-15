package main

import "testing"

func TestPatterns(t *testing.T) {
	patterns := []string{"a+b", "(a+b)", "(a+b)*", "(a+b)+(c+d)", "a+c+d+e+f+g", "(a+b)*+c", "(a+b)*+(a+bb)", "(a+b)*+((a+bb))", "((a+b*)*abb)*"}

	for i, pattern := range patterns {
		nfa := buildNFA(pattern)
		if nfa == nil {
			t.Errorf("Test %d/%d failed for pattern %s", i, len(patterns), pattern)
		}
	}
}
