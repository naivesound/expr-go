package expr

import "testing"

func BenchmarkExprParse(b *testing.B) {
	env := map[string]Var{}
	funcs := map[string]Func{
		"plusone": func(c *FuncContext) Num {
			return c.Args[0].Eval() + 1
		},
	}
	expr := "x=2+3*(x/(42+plusone(x))),x"
	for i := 0; i < b.N; i++ {
		Parse(expr, env, funcs)
	}
}

func BenchmarkExprEval(b *testing.B) {
	env := map[string]Var{}
	funcs := map[string]Func{
		"plusone": func(c *FuncContext) Num {
			return c.Args[0].Eval() + 1
		},
	}
	expr := "x=2+3*(x/(42+plusone(x))),x"
	if e, err := Parse(expr, env, funcs); err != nil {
		b.Error(err)
	} else {
		for i := 0; i < b.N; i++ {
			e.Eval()
		}
	}
}
