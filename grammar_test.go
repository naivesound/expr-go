package expr

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	// Let's preter there is no "&" operator, but there is "&&"
	defer func() {
		ops["&"] = bitwiseAnd
	}()
	delete(ops, "&")

	for s, parts := range map[string][]string{
		"2":         {"2"},
		"2+3/234.0": {"2", "+", "3", "/", "234.0"},
		"2+-3":      {"2", "+", "-u", "3"},
		"2--3":      {"2", "-", "-u", "3"},
		"-(-2)":     {"-u", "(", "-u", "2", ")"},
		"foo":       {"foo"},
		"1>2":       {"1", ">", "2"},
		"1>-2":      {"1", ">", "-u", "2"},
		"1>>2":      {"1", ">>", "2"},
		"1>>-2":     {"1", ">>", "-u", "2"},
		"1>>!2":     {"1", ">>", "!", "2"},
		"1&&2":      {"1", "&&", "2"},
		"1&&":       {"1", "&&"},
		"1&&&":      nil, // This should return an error: 'no such operator &'
	} {
		if tokens, err := tokenize([]rune(s)); err != nil {
			if parts != nil {
				t.Error(err)
			}
		} else if len(tokens) != len(parts) {
			t.Error(tokens, parts)
		} else {
			for i, tok := range tokens {
				if tok != parts[i] {
					t.Error(tokens, parts)
					break
				}
			}
		}
	}
}

func TestParse(t *testing.T) {
	env := map[string]Var{
		"x": NewVarExpr(5),
	}
	funcs := map[string]Func{
		"add3": NewFunc(func(args FuncArgs, env FuncEnv) Num {
			return args[0].Eval() + args[1].Eval() + args[2].Eval()
		}),
	}
	for input, result := range map[string]Num{
		"":  0,
		"2": 2,
		"x": 5,

		"-2": -2,
		"~2": -3,
		"!2": 0,
		"!0": 1,

		"3+2":       5,
		"3/2":       1.5,
		"2+3/2":     2 + 3/2.0,
		"4/2+8*4/2": 18,

		"2*x": 10,
		"2/x": 2 / 5.0,

		"4*(2+8)+4/2": 42,

		"2, 3, 5":  5,
		"2+3, 5*3": 15,

		"z=10":     10,
		"y=10,x+y": 15,

		"2+add3(3, 7, 9)":             21,
		"2+add3(3, add3(1, 2, 3), 9)": 20,
	} {
		if e, err := Parse(input, env, funcs); err != nil {
			t.Error(err)
		} else if n := e.Eval(); n != result {
			t.Error(n, result)
		}
	}
}

func TestExprString(t *testing.T) {
	env := map[string]Var{
		"x": NewVarExpr(5),
	}
	funcs := map[string]Func{
		"plusone": NewFunc(func(args FuncArgs, env FuncEnv) Num {
			return args[0].Eval() + 1
		}),
	}
	e, _ := Parse("-2+plusone(x)", env, funcs)
	if s := fmt.Sprintf("%v", e); s != "fn[<8>(<1>(#2), fn[{5}])]" {
		t.Error(e)
	}
}
