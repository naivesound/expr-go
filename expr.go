package expr

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

type Num float64

var (
	ErrBadCall         = errors.New("function call expected")
	ErrBadAssignment   = errors.New("variable expected in assignment")
	ErrBadOp           = errors.New("unknown operator or function")
	ErrOperandMissing  = errors.New("missing operand")
	ErrOperatorMissing = errors.New("missing operator")
	ErrParen           = errors.New("parenthesis mismatch")
)

// Supported arithmetic operations
type arithOp int

const (
	unaryMinus arithOp = iota + 1
	unaryLogicalNot
	unaryBitwiseNot
	unarySqrt

	multiply
	divide
	remainder

	plus
	minus

	shl
	shr

	lessThan
	lessOrEquals
	greaterThan
	greaterOrEquals
	equals
	notEquals

	bitwiseAnd
	bitwiseXor
	bitwiseOr

	logicalAnd
	logicalOr

	assign
)

var ops = map[string]arithOp{
	"-u": unaryMinus, "!": unaryLogicalNot, "~": unaryBitwiseNot, "√": unarySqrt,
	"*": multiply, "/": divide, "%": remainder,
	"+": plus, "-": minus,
	"<<": shl, ">>": shr,
	"<": lessThan, "<=": lessOrEquals, ">": greaterThan, ">=": greaterOrEquals,
	"==": equals, "!=": notEquals,
	"&": bitwiseAnd, "^": bitwiseXor, "|": bitwiseOr,
	"&&": logicalAnd, "||": logicalOr,
	"=": assign,
}

func isUnary(op arithOp) bool {
	return op >= unaryMinus && op <= unarySqrt
}
func isLeftAssoc(op arithOp) bool {
	return !isUnary(op) && op != assign
}
func boolNum(b bool) Num {
	if b {
		return 1
	} else {
		return 0
	}
}

type Expr interface {
	Eval() Num
}

// Constant expression always returns the same value when evaluated
type constExpr struct {
	value Num
}

func (e *constExpr) Eval() Num {
	return e.value
}

func (e *constExpr) String() string {
	return fmt.Sprintf("#%v", e.value)
}

// Mutable variable expression returns the currently stored value of the variable
type Var interface {
	Expr
	Set(value Num)
}
type varExpr struct {
	value Num
}

func NewVarExpr(value Num) Var {
	return &varExpr{value: value}
}
func (e *varExpr) Eval() Num {
	return e.value
}
func (e *varExpr) Set(value Num) {
	e.value = value
}
func (e *varExpr) String() string {
	return fmt.Sprintf("{%v}", e.value)
}

// Function expression returns the result of the function
type Func interface {
	Bind(args []Expr) Expr
}

type FuncEnv map[string]Num
type FuncArgs []Expr

type EvalFunc func(FuncArgs, FuncEnv) Num
type funcBinder struct {
	eval EvalFunc
}

type simpleFunc struct {
	eval EvalFunc
	env  FuncEnv
	args FuncArgs
}

func NewFunc(eval EvalFunc) Func {
	return &funcBinder{eval: eval}
}

func (f *funcBinder) Bind(args []Expr) Expr {
	env := map[string]Num{}
	return &simpleFunc{eval: f.eval, args: args, env: env}
}

func (e *simpleFunc) Eval() Num {
	return e.eval(e.args, e.env)
}

func (e *simpleFunc) String() string {
	return fmt.Sprintf("fn%v", e.args)
}

var lastArgFunc = NewFunc(func(args FuncArgs, env FuncEnv) Num {
	var result Num
	for _, arg := range args {
		result = arg.Eval()
	}
	return result
})

// Operator expression returns the result of the operator applied to 1 or 2 arguments
type unaryExpr struct {
	op  arithOp
	arg Expr
}

func newUnaryExpr(op arithOp, arg Expr) Expr {
	return &unaryExpr{op: op, arg: arg}
}
func (e *unaryExpr) Eval() (res Num) {
	switch e.op {
	case unaryMinus:
		res = -e.arg.Eval()
	case unaryBitwiseNot:
		// Bitwise operation can only be applied to integer values
		res = Num(^int64(e.arg.Eval()))
	case unaryLogicalNot:
		res = boolNum(e.arg.Eval() == 0)
	case unarySqrt:
		res = Num(math.Sqrt(float64(e.arg.Eval())))
	}
	return res
}
func (e *unaryExpr) String() string {
	return fmt.Sprintf("<%v>(%v)", e.op, e.arg)
}

type binaryExpr struct {
	op arithOp
	a  Expr
	b  Expr
}

func newBinaryExpr(op arithOp, a, b Expr) (Expr, error) {
	if op == assign {
		if _, ok := a.(*varExpr); !ok {
			return nil, ErrBadAssignment
		}
	}
	return &binaryExpr{op: op, a: a, b: b}, nil
}

func (e *binaryExpr) Eval() (res Num) {
	switch e.op {
	case multiply:
		res = e.a.Eval() * e.b.Eval()
	case divide:
		tmp := e.b.Eval()
		if tmp != 0 {
			res = e.a.Eval() / tmp
		}
	case remainder:
		tmp := e.b.Eval()
		if tmp != 0 {
			res = Num(math.Remainder(float64(e.a.Eval()), float64(tmp)))
		}
	case plus:
		res = e.a.Eval() + e.b.Eval()
	case minus:
		res = e.a.Eval() - e.b.Eval()
	case shl:
		res = Num(int64(e.a.Eval()) << uint(e.b.Eval()))
	case shr:
		res = Num(int64(e.a.Eval()) >> uint(e.b.Eval()))
	case lessThan:
		res = boolNum(e.a.Eval() < e.b.Eval())
	case lessOrEquals:
		res = boolNum(e.a.Eval() <= e.b.Eval())
	case greaterThan:
		res = boolNum(e.a.Eval() > e.b.Eval())
	case greaterOrEquals:
		res = boolNum(e.a.Eval() >= e.b.Eval())
	case equals:
		res = boolNum(e.a.Eval() == e.b.Eval())
	case notEquals:
		res = boolNum(e.a.Eval() != e.b.Eval())
	case bitwiseAnd:
		return Num(int64(e.a.Eval()) & int64(e.b.Eval()))
	case bitwiseXor:
		return Num(int64(e.a.Eval()) ^ int64(e.b.Eval()))
	case bitwiseOr:
		return Num(int64(e.a.Eval()) | int64(e.b.Eval()))
	case logicalAnd:
		res = boolNum(e.a.Eval() != 0 && e.b.Eval() != 0)
	case logicalOr:
		res = boolNum(e.a.Eval() != 0 || e.b.Eval() != 0)
	case assign:
		res = e.b.Eval()
		e.a.(*varExpr).Set(res)
	}
	return res
}

func (e *binaryExpr) String() string {
	return fmt.Sprintf("<%v>(%v, %v)", e.op, e.a, e.b)
}

func tokenize(input []rune) (tokens []string, err error) {
	pos := 0
	expectVal := true
	for pos < len(input) {
		tok := []rune{}
		c := input[pos]
		if unicode.IsSpace(c) {
			pos++
			continue
		}
		if unicode.IsNumber(c) {
			if !expectVal {
				return nil, ErrOperandMissing
			}
			expectVal = false
			for (c == '.' || unicode.IsNumber(c)) && pos < len(input) {
				tok = append(tok, input[pos])
				pos++
				if pos < len(input) {
					c = input[pos]
				} else {
					c = 0
				}
			}
		} else if c == '-' {
			// Minus, unary or binary
			if expectVal {
				expectVal = true
				tok = append(tok, '-', 'u')
			} else {
				expectVal = true
				tok = append(tok, '-')
			}
			pos++
		} else if unicode.IsLetter(c) {
			if !expectVal {
				return nil, ErrOperandMissing
			}
			expectVal = false
			for unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' && pos < len(input) {
				tok = append(tok, input[pos])
				pos++
				if pos < len(input) {
					c = input[pos]
				} else {
					c = 0
				}
			}
		} else if c == '(' || c == ')' || c == ',' {
			expectVal = (c == '(' || c == ',')
			tok = append(tok, c)
			pos++
		} else {
			expectVal = true
			var lastOp string
			for !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsSpace(c) &&
				c != '_' && c != '(' && c != ')' && c != ',' && c != '-' && pos < len(input) {
				if _, ok := ops[string(tok)+string(input[pos])]; ok {
					tok = append(tok, input[pos])
					lastOp = string(tok)
				} else if lastOp == "" {
					tok = append(tok, input[pos])
				} else {
					break
				}
				pos++
				if pos < len(input) {
					c = input[pos]
				} else {
					c = 0
				}
			}
			if lastOp == "" {
				return nil, ErrBadOp
			}
		}
		tokens = append(tokens, string(tok))
	}
	return tokens, nil
}

// Simple string stack implementation
type stringStack []string

func (ss *stringStack) Push(s string) {
	*ss = append(*ss, s)
}
func (ss *stringStack) Peek() string {
	if l := len(*ss); l == 0 {
		return ""
	} else {
		return (*ss)[l-1]
	}
}
func (ss *stringStack) Pop() string {
	if l := len(*ss); l == 0 {
		return ""
	} else {
		s := (*ss)[l-1]
		*ss = (*ss)[:l-1]
		return s
	}
}

// Simple expression stack implementation
type exprStack []Expr

func (es *exprStack) Push(e Expr) {
	*es = append(*es, e)
}
func (es *exprStack) Peek() Expr {
	if l := len(*es); l == 0 {
		return nil
	} else {
		return (*es)[l-1]
	}
}
func (es *exprStack) Pop() Expr {
	if l := len(*es); l == 0 {
		return nil
	} else {
		e := (*es)[l-1]
		*es = (*es)[:l-1]
		return e
	}
}

func Parse(input string, vars map[string]Var, funcs map[string]Func) (Expr, error) {
	os := stringStack{}
	es := exprStack{}

	expectArgs := false

	input = "(" + input + ")" // To allow multiple return values
	if tokens, err := tokenize([]rune(input)); err != nil {
		return nil, err
	} else {
		for _, token := range tokens {
			if token == "(" {
				expectArgs = false
				os.Push("(")
				es.Push(nil)
			} else if expectArgs {
				return nil, ErrBadCall
			} else if token == ")" {
				for len(os) > 0 && os.Peek() != "(" {
					if expr, err := bind(os.Pop(), funcs, &es); err != nil {
						return nil, err
					} else {
						es.Push(expr)
					}
				}
				if len(os) == 0 {
					return nil, ErrParen
				}
				os.Pop()
				if len(os) > 0 && funcs[os.Peek()] != nil {
					if expr, err := bind(os.Pop(), funcs, &es); err != nil {
						// Should never happen XXX why?
						return nil, err
					} else {
						es.Push(expr)
					}
				} else {
					es.Push(lastArgFunc.Bind(parseArgs(&es)))
				}
			} else if n, err := strconv.ParseFloat(token, 64); err == nil {
				// Number
				es.Push(&constExpr{value: Num(n)})
			} else if _, ok := funcs[token]; ok {
				// Function
				os.Push(token)
				expectArgs = true
			} else if token == "," {
				for len(os) > 0 && os.Peek() != "(" {
					if expr, err := bind(os.Pop(), funcs, &es); err != nil {
						return nil, err
					} else {
						es.Push(expr)
					}
				}
				if len(os) == 0 {
					// Should never happen as long as we wrap input string in extra parenthesis
					return nil, ErrParen
				}
			} else if op, ok := ops[token]; ok {
				o2 := os.Peek()
				for ops[o2] != 0 && ((isLeftAssoc(op) && op >= ops[o2]) || op > ops[o2]) {
					if expr, err := bind(o2, funcs, &es); err != nil {
						return nil, err
					} else {
						es.Push(expr)
					}
					os.Pop()
					o2 = os.Peek()
				}
				os.Push(token)
			} else {
				// Variable
				if v, ok := vars[token]; ok {
					es.Push(v)
				} else {
					v = NewVarExpr(0)
					vars[token] = v
					es.Push(v)
				}
			}
		}
		for len(os) > 0 {
			op := os.Pop()
			if op == "(" || op == ")" {
				return nil, ErrParen
			}
			if expr, err := bind(op, funcs, &es); err != nil {
				return nil, err
			} else {
				es.Push(expr)
			}
		}
		if len(es) == 0 {
			// Should never happen as long as we wrap input in extra parenthesis
			return &constExpr{}, nil
		} else {
			e := es.Pop()
			return e, nil
		}
	}
}

func bind(name string, funcs map[string]Func, stack *exprStack) (Expr, error) {
	if f, ok := funcs[name]; ok {
		return f.Bind(parseArgs(stack)), nil
	} else if op, ok := ops[name]; ok {
		if isUnary(op) {
			if stack.Peek() == nil {
				return nil, ErrOperandMissing
			} else {
				return newUnaryExpr(op, stack.Pop()), nil
			}
		} else {
			b := stack.Pop()
			a := stack.Pop()
			if a == nil || b == nil {
				return nil, ErrOperandMissing
			}
			return newBinaryExpr(op, a, b)
		}
	} else {
		// Should never happen, bad operators are filtered in tokenizer
		return nil, ErrBadOp
	}
}

func parseArgs(stack *exprStack) []Expr {
	args := []Expr{}
	for len(*stack) > 0 && stack.Peek() != nil {
		args = append([]Expr{stack.Pop()}, args...)
	}
	if len(*stack) > 0 {
		stack.Pop()
	}
	return args
}
