# expr - tiny yet super fast math expression evaluator in Go

[![GoDoc](https://godoc.org/github.com/naivesound/expr?status.svg)](https://godoc.org/github.com/naivesound/expr)
[![Build Status](https://travis-ci.org/naivesound/expr.svg?branch=master)](https://travis-ci.org/naivesound/expr)
[![codecov.io](https://codecov.io/github/naivesound/expr/coverage.svg?branch=master)](https://codecov.io/github/naivesound/expr?branch=master)

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
