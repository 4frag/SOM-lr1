package tasks

import (
	"fmt"
	"math"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type ExpressionParser struct {
	program *vm.Program
}


func NewExpressionParser(expression string) (*ExpressionParser, error) {
	env := map[string]interface{}{
		"x": 0.0,
		"y": 0.0,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"exp":   math.Exp,
		"log":   math.Log,
		"log10": math.Log10,
		"sqrt":  math.Sqrt,
		"abs":   math.Abs,
		"pow":   math.Pow,
		"pi":    math.Pi,
		"e":     math.E,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("ошибка компиляции выражения: %v", err)
	}

	return &ExpressionParser{program: program}, nil
}

func (p *ExpressionParser) Eval(x, y float64) (float64, error) {
	env := map[string]interface{}{
		"x": x,
		"y": y,
		"sin":   math.Sin,
		"cos":   math.Cos,
		"tan":   math.Tan,
		"exp":   math.Exp,
		"log":   math.Log,
		"log10": math.Log10,
		"sqrt":  math.Sqrt,
		"abs":   math.Abs,
		"pow":   math.Pow,
		"pi":    math.Pi,
		"e":     math.E,
	}

	result, err := expr.Run(p.program, env)
	if err != nil {
		return 0, fmt.Errorf("ошибка вычисления: %v", err)
	}

	switch v := result.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("неожиданный тип результата: %T", result)
	}
}

func CreateFunction(expression string) (func(x, y float64) float64, error) {
	parser, err := NewExpressionParser(expression)
	if err != nil {
		return nil, err
	}

	return func(x, y float64) float64 {
		result, err := parser.Eval(x, y)
		if err != nil {
			return 0.0
		}
		return result
	}, nil
}