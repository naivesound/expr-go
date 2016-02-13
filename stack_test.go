package expr

import "testing"

func TestStringStack(t *testing.T) {
	s := stringStack{}
	s.Push("foo")
	s.Push("bar")
	s.Push("baz")
	if len(s) != 3 {
		t.Error()
	}
	if s.Peek() != "baz" {
		t.Error()
	}
	if s.Pop() != "baz" {
		t.Error()
	}
	if len(s) != 2 {
		t.Error()
	}
	if s.Pop() != "bar" {
		t.Error()
	}
	if s.Pop() != "foo" {
		t.Error()
	}
	if s.Pop() != "" || s.Peek() != "" {
		t.Error()
	}
	if s.Pop() != "" {
		t.Error()
	}
	if len(s) != 0 {
		t.Error()
	}
}

func TestExprStack(t *testing.T) {
	s := exprStack{}
	s.Push(&constExpr{value: 1})
	s.Push(&constExpr{value: 2})
	s.Push(&constExpr{value: 3})
	if len(s) != 3 {
		t.Error()
	}
	if s.Peek().Eval() != 3 {
		t.Error()
	}
	if s.Pop().Eval() != 3 {
		t.Error()
	}
	if len(s) != 2 {
		t.Error()
	}
	s.Pop()
	s.Pop()
	if s.Pop() != nil || s.Peek() != nil {
		t.Error()
	}
	if s.Pop() != nil {
		t.Error()
	}
	if len(s) != 0 {
		t.Error()
	}
}
