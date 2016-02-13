package expr

import "testing"

func TestConstExpr(t *testing.T) {
	e := &constExpr{value: 3}
	if n := e.Eval(); n != 3 {
		t.Error(n)
	}
}

func TestVarExpr(t *testing.T) {
	e := NewVarExpr(3)
	if n := e.Eval(); n != 3 {
		t.Error(n)
	}
	e.Set(5)
	if n := e.Eval(); n != 5 {
		t.Error(n)
	}
}

func TestFuncExpr(t *testing.T) {
	f := NewFunc(func(args FuncArgs, env FuncEnv) Num {
		env["accum"] = env["accum"] + args[0].Eval()
		return env["accum"]
	})
	two := &constExpr{value: 2}
	x := NewVarExpr(0)
	sum := f.Bind([]Expr{two})
	sumvar := f.Bind([]Expr{x})

	if n := sum.Eval(); n != 2 {
		t.Error(n)
	}
	if n := sum.Eval(); n != 4 {
		t.Error(n)
	}

	if n := sumvar.Eval(); n != 0 {
		t.Error(n)
	}
	x.Set(2)
	if n := sumvar.Eval(); n != 2 {
		t.Error(n)
	}
	x.Set(5)
	if n := sumvar.Eval(); n != 7 {
		t.Error(n)
	}
	x.Set(8)
	if n := sumvar.Eval(); n != 15 {
		t.Error(n)
	}
}

func TestLastArgFunc(t *testing.T) {
	args := []Expr{
		&constExpr{value: 2},
		NewVarExpr(3),
		NewVarExpr(7),
	}
	f := lastArgFunc.Bind(args)
	if n := f.Eval(); n != 7 {
		t.Error(n)
	}
}

func TestUnaryExpr(t *testing.T) {
	for e, res := range map[Expr]Num{
		newUnaryExpr(unaryMinus, &constExpr{5}):      -5,
		newUnaryExpr(unarySqrt, &constExpr{9}):       3,
		newUnaryExpr(unaryBitwiseNot, &constExpr{9}): -10,
		newUnaryExpr(unaryLogicalNot, &constExpr{9}): 0,
		newUnaryExpr(unaryLogicalNot, &constExpr{0}): 1,
	} {
		if n := e.Eval(); n != res {
			t.Error(n, res)
		}
	}
}

func TestBinaryExpr(t *testing.T) {
	for e, res := range map[Expr]Num{
		newBinaryExpr(multiply, &constExpr{9}, &constExpr{4}):  36,
		newBinaryExpr(divide, &constExpr{9}, &constExpr{4}):    9 / 4.0,
		newBinaryExpr(remainder, &constExpr{9}, &constExpr{4}): 1,

		newBinaryExpr(plus, &constExpr{5}, &constExpr{3}):  8,
		newBinaryExpr(minus, &constExpr{9}, &constExpr{4}): 5,

		newBinaryExpr(shl, &constExpr{5}, &constExpr{1}): 10,
		newBinaryExpr(shr, &constExpr{9}, &constExpr{1}): 4,

		newBinaryExpr(lessThan, &constExpr{5}, &constExpr{5}):        0,
		newBinaryExpr(lessOrEquals, &constExpr{9}, &constExpr{9}):    1,
		newBinaryExpr(greaterThan, &constExpr{5}, &constExpr{3}):     1,
		newBinaryExpr(greaterOrEquals, &constExpr{9}, &constExpr{4}): 1,
		newBinaryExpr(equals, &constExpr{5}, &constExpr{3}):          0,
		newBinaryExpr(equals, &constExpr{5}, NewVarExpr(5)):          1,
		newBinaryExpr(notEquals, &constExpr{9}, &constExpr{0}):       1,
		newBinaryExpr(notEquals, &constExpr{5}, NewVarExpr(5)):       0,

		newBinaryExpr(bitwiseAnd, &constExpr{10}, &constExpr{7}): 2,
		newBinaryExpr(bitwiseOr, &constExpr{9}, &constExpr{4}):   13,
		newBinaryExpr(bitwiseXor, &constExpr{9}, &constExpr{2}):  11,

		newBinaryExpr(logicalAnd, &constExpr{9}, &constExpr{4}): 1,
		newBinaryExpr(logicalOr, &constExpr{9}, &constExpr{4}):  1,

		newBinaryExpr(assign, NewVarExpr(0), &constExpr{4}): 4,
	} {
		if n := e.Eval(); n != res {
			t.Error(n, res)
		}
	}
}
