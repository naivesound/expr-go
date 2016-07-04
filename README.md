# expr - tiny yet super fast math expression evaluator in Go

[![GoDoc](https://godoc.org/github.com/naivesound/expr-go?status.svg)](https://godoc.org/github.com/naivesound/expr-go)
[![Build Status](https://travis-ci.org/naivesound/expr-go.svg?branch=master)](https://travis-ci.org/naivesound/expr-go)
[![codecov.io](https://codecov.io/github/naivesound/expr-go/coverage.svg?branch=master)](https://codecov.io/github/naivesound/expr-go?branch=master)

## Example

```go
vars := map[string]expr.Var{
	"x": NewVar(5),
}
funcs := map[string]expr.Func{
	"next": NewFunc(func(args expr.FuncArgs, env expr.FuncEnv) expr.Num {
		return expr.Num(args[0].Eval()+1)
	}),
}
vars["x"].Set(5)
e, err := expr.Parse("y=x+5/next(x)", vars, funcs)
if err != nil {
	log.Fatal(err)
}
result := e.Eval()
log.Println(result)
log.Println(vars["y"])
```

## Performance

The goal is to speed up frequent evaluations of immutable expressions.
For example, user enters a formula to calculate some metrics and backend does
the calculation. Formula itself changes rarely, but is evaluated very often.

On my dual-core laptop the benchmark shows that:

- Parsing expression `x=2+3*(x/(42+plusone(x))),x` takes 14us, evaluation takes **65ns**.
- Parsing a more complex expression of 280 chars takes 112us, evaluation takes **649ns**.
- Parsing an even more complex expresison of ~3k chars takes 1ms, evaluation takes **7us**.

This probably means that time grows linearly depending on the formula length
and the performance is really good for short (human-generated) formulas.
