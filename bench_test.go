package expr

import "testing"

var expr = "x=2+3*(x/(42+plusone(x))),x"

func bench(n int, eval bool, b *testing.B) {
	s := "0"
	for i := 0; i < n; i++ {
		s = s + "," + expr
	}
	env := map[string]Var{}
	funcs := map[string]Func{
		"plusone": func(c *FuncContext) Num {
			return c.Args[0].Eval() + 1
		},
	}
	if eval {
		if e, err := Parse(s, env, funcs); err != nil {
			b.Error(err)
		} else if eval {
			for i := 0; i < b.N; i++ {
				e.Eval()
			}
		}
	} else {
		for i := 0; i < b.N; i++ {
			if _, err := Parse(s, env, funcs); err != nil {
				b.Error(err)
			}
		}
	}
}

func BenchmarkExprParse(b *testing.B) {
	bench(1, false, b)
}

func BenchmarkExprParse10(b *testing.B) {
	bench(10, false, b)
}

func BenchmarkExprParse100(b *testing.B) {
	bench(100, false, b)
}

func BenchmarkExprEval(b *testing.B) {
	bench(1, true, b)
}

func BenchmarkExprEval10(b *testing.B) {
	bench(10, true, b)
}
func BenchmarkExprEval100(b *testing.B) {
	bench(100, true, b)
}
