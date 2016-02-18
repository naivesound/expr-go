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
	ErrParen                = errors.New("parenthesis mismatch")
	ErrUnexpectedNumber     = errors.New("unexpected number")
	ErrUnexpectedIdentifier = errors.New("unexpected identifier")

	ErrBadCall         = errors.New("function call expected")
	ErrBadAssignment   = errors.New("variable expected in assignment")
	ErrBadOp           = errors.New("unknown operator or function")
	ErrOperandMissing  = errors.New("missing operand")
	ErrOperatorMissing = errors.New("missing operator")
)

// Supported arithmetic operations
type arithOp int

const (
	unaryMinus arithOp = iota + 1
	unaryLogicalNot
	unaryBitwiseNot

	power
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
	comma
)

var ops = map[string]arithOp{
	"-u": unaryMinus, "!u": unaryLogicalNot, "^u": unaryBitwiseNot,
	"**": power, "*": multiply, "/": divide, "%": remainder,
	"+": plus, "-": minus,
	"<<": shl, ">>": shr,
	"<": lessThan, "<=": lessOrEquals, ">": greaterThan, ">=": greaterOrEquals,
	"==": equals, "!=": notEquals,
	"&": bitwiseAnd, "^": bitwiseXor, "|": bitwiseOr,
	"&&": logicalAnd, "||": logicalOr,
	"=": assign, ",": comma,
}

func isUnary(op arithOp) bool {
	return op >= unaryMinus && op <= unaryBitwiseNot
}
func isLeftAssoc(op arithOp) bool {
	return !isUnary(op) && op != assign && op != power && op != comma
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
	Get() Num
}
type varExpr struct {
	value Num
}

func NewVar(value Num) Var {
	return &varExpr{value: value}
}
func (e *varExpr) Eval() Num {
	return e.value
}
func (e *varExpr) Set(value Num) {
	e.value = value
}
func (e *varExpr) Get() Num {
	return e.value
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
	case power:
		res = Num(math.Pow(float64(e.a.Eval()), float64(e.b.Eval())))
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
		if a := e.a.Eval(); a != 0 {
			if b := e.b.Eval(); b != 0 {
				res = b
			}
		}
	case logicalOr:
		if a := e.a.Eval(); a != 0 {
			res = a
		} else if b := e.b.Eval(); b != 0 {
			res = b
		}
	case assign:
		res = e.b.Eval()
		e.a.(*varExpr).Set(res)
	case comma:
		e.a.Eval()
		res = e.b.Eval()
	}
	return res
}

func (e *binaryExpr) String() string {
	return fmt.Sprintf("<%v>(%v, %v)", e.op, e.a, e.b)
}

const (
	tokNumber = 1 << iota
	tokWord
	tokOp
	tokOpen
	tokClose
)

func tokenize(input []rune) (tokens []string, err error) {
	pos := 0
	expected := tokOpen | tokNumber | tokWord
	for pos < len(input) {
		tok := []rune{}
		c := input[pos]
		if unicode.IsSpace(c) {
			pos++
			continue
		}
		if unicode.IsNumber(c) {
			if expected&tokNumber == 0 {
				return nil, ErrUnexpectedNumber
			}
			expected = tokOp | tokClose
			for (c == '.' || unicode.IsNumber(c)) && pos < len(input) {
				tok = append(tok, input[pos])
				pos++
				if pos < len(input) {
					c = input[pos]
				} else {
					c = 0
				}
			}
		} else if unicode.IsLetter(c) {
			if expected&tokWord == 0 {
				return nil, ErrUnexpectedIdentifier
			}
			expected = tokOp | tokOpen | tokClose
			for unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' && pos < len(input) {
				tok = append(tok, input[pos])
				pos++
				if pos < len(input) {
					c = input[pos]
				} else {
					c = 0
				}
			}
		} else if c == '(' || c == ')' {
			tok = append(tok, c)
			pos++
			if c == '(' && (expected&tokOpen) != 0 {
				expected = tokNumber | tokWord | tokOpen | tokClose
			} else if c == ')' && (expected&tokClose) != 0 {
				expected = tokOp | tokClose
			} else {
				return nil, ErrParen
			}
		} else if c == '-' || c == '^' || c == '!' {
			// Minus, unary or binary
			if expected&tokOp == 0 {
				tok = append(tok, c, 'u')
				expected = tokNumber | tokWord | tokOpen
			} else {
				tok = append(tok, c)
				expected = tokNumber | tokWord | tokOpen
			}
			pos++
		} else {
			if expected&tokOp == 0 {
				return nil, ErrOperandMissing
			}
			expected = tokNumber | tokWord | tokOpen
			var lastOp string
			for !unicode.IsLetter(c) && !unicode.IsNumber(c) && !unicode.IsSpace(c) &&
				c != '_' && c != '(' && c != ')' && c != '^' && c != '-' && pos < len(input) {
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

	if tokens, err := tokenize([]rune(input)); err != nil {
		return nil, err
	} else {
		for _, token := range tokens {
			if token == "(" {
				expectArgs = false
				os.Push("(")
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
				os.Pop() // remove open paren
				if len(os) > 0 {
					if expr, err := bind(os.Pop(), funcs, &es); err != nil {
						return nil, err
					} else {
						es.Push(expr)
					}
				}
			} else if n, err := strconv.ParseFloat(token, 64); err == nil {
				// Number
				es.Push(&constExpr{value: Num(n)})
			} else if _, ok := funcs[token]; ok {
				// Function
				os.Push(token)
				expectArgs = true
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
					v = NewVar(0)
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
			return &constExpr{}, nil
		} else {
			e := es.Pop()
			return e, nil
		}
	}
}

func bind(name string, funcs map[string]Func, stack *exprStack) (Expr, error) {
	if f, ok := funcs[name]; ok {
		args := []Expr{}
		if len(*stack) > 0 {
			if stack.Peek() != nil {
				args = list(stack.Pop())
			}
		} else {
			return nil, ErrBadCall
		}
		return f.Bind(args), nil
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
		return nil, ErrBadOp
	}
}

func list(e Expr) []Expr {
	if b, ok := e.(*binaryExpr); ok && b.op == comma {
		return append([]Expr{b.a}, list(b.b)...)
	} else {
		return []Expr{e}
	}
}
