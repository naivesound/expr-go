package expr

import "testing"

func TestConstExpr(t *testing.T) {
	e := &constExpr{value: 3}
	if n := e.Eval(); n != 3 {
		t.Error(n)
	}
}

func TestVar(t *testing.T) {
	e := NewVar(3)
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
	x := NewVar(0)
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
		NewVar(3),
		NewVar(7),
	}
	f := lastArgFunc.Bind(args)
	if n := f.Eval(); n != 7 {
		t.Error(n)
	}
}

func TestUnaryExpr(t *testing.T) {
	for e, res := range map[Expr]Num{
		newUnaryExpr(unaryMinus, &constExpr{5}):      -5,
		newUnaryExpr(unaryBitwiseNot, &constExpr{9}): -10,
		newUnaryExpr(unaryLogicalNot, &constExpr{9}): 0,
		newUnaryExpr(unaryLogicalNot, &constExpr{0}): 1,
	} {
		if n := e.Eval(); n != res {
			t.Error(e, n, res)
		}
	}
}

func TestBinaryExpr(t *testing.T) {
	for e, res := range map[Expr]Num{
		&binaryExpr{power, &constExpr{9}, &constExpr{4}}:      6561,
		&binaryExpr{multiply, &constExpr{9}, &constExpr{4}}:   36,
		&binaryExpr{divide, &constExpr{9}, &constExpr{4}}:     9.0 / 4.0,
		&binaryExpr{remainder, &constExpr{9}, &constExpr{4}}:  1,
		&binaryExpr{remainder, &constExpr{9}, &constExpr{9}}:  0,
		&binaryExpr{remainder, &constExpr{9}, &constExpr{0}}:  0,
		&binaryExpr{remainder, &constExpr{-9}, &constExpr{9}}: 0,
		&binaryExpr{remainder, &constExpr{-9}, &constExpr{8}}: -1,

		&binaryExpr{plus, &constExpr{5}, &constExpr{3}}:  8,
		&binaryExpr{minus, &constExpr{9}, &constExpr{4}}: 5,

		&binaryExpr{shl, &constExpr{5}, &constExpr{1}}: 10,
		&binaryExpr{shr, &constExpr{9}, &constExpr{1}}: 4,

		&binaryExpr{lessThan, &constExpr{5}, &constExpr{5}}:        0,
		&binaryExpr{lessOrEquals, &constExpr{9}, &constExpr{9}}:    1,
		&binaryExpr{greaterThan, &constExpr{5}, &constExpr{3}}:     1,
		&binaryExpr{greaterOrEquals, &constExpr{9}, &constExpr{4}}: 1,
		&binaryExpr{equals, &constExpr{5}, &constExpr{3}}:          0,
		&binaryExpr{equals, &constExpr{5}, NewVar(5)}:              1,
		&binaryExpr{notEquals, &constExpr{9}, &constExpr{0}}:       1,
		&binaryExpr{notEquals, &constExpr{5}, NewVar(5)}:           0,

		&binaryExpr{bitwiseAnd, &constExpr{10}, &constExpr{7}}: 2,
		&binaryExpr{bitwiseOr, &constExpr{9}, &constExpr{4}}:   13,
		&binaryExpr{bitwiseXor, &constExpr{9}, &constExpr{2}}:  11,

		// Returns last argument if true, or 0 if false
		&binaryExpr{logicalAnd, &constExpr{9}, &constExpr{4}}: 4,
		&binaryExpr{logicalAnd, &constExpr{9}, &constExpr{0}}: 0,
		// Returns first argument if true, or second if false
		&binaryExpr{logicalOr, &constExpr{3}, &constExpr{4}}: 3,
		&binaryExpr{logicalOr, &constExpr{0}, &constExpr{4}}: 4,
		&binaryExpr{logicalOr, &constExpr{0}, &constExpr{0}}: 0,

		&binaryExpr{assign, NewVar(0), &constExpr{4}}: 4,
	} {
		if n := e.Eval(); n != res {
			t.Error(e, n, res)
		}
	}
}
