//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Evaluator struct {
	m     map[string][]string
	stack []int
}

// declare operations and their balance changing
var operations = map[string][]int{
	"+":    {2, -1},
	"-":    {2, -1},
	"*":    {2, -1},
	"/":    {2, -1},
	"dup":  {1, 1},
	"drop": {1, -1},
	"over": {2, 1},
	"swap": {2, 0},
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	eval := Evaluator{
		m:     make(map[string][]string),
		stack: make([]int, 0, 7),
	}
	for str := range operations {
		eval.m[str] = append(eval.m[str], str)
	}
	return &eval
}

func (e *Evaluator) doCommand(arr []string) error {

	//checking ability to do this operations
	balance := len(e.stack)
	for _, str := range arr {
		if _, err := strconv.Atoi(str); err == nil {
			balance++
			continue
		}
		for key, value := range operations {
			if key == str {
				if balance < value[0] {
					return errors.New("not enought args on stack")
				}
				balance += value[1]
			}
		}
	}

	tmp := e.stack

	for _, str := range arr {
		switch str {
		case "+":
			sum := e.stack[len(e.stack)-1] + e.stack[len(e.stack)-2]
			e.stack = e.stack[:len(e.stack)-1]
			e.stack[len(e.stack)-1] = sum
		case "-":
			diff := e.stack[len(e.stack)-2] - e.stack[len(e.stack)-1]
			e.stack = e.stack[:len(e.stack)-1]
			e.stack[len(e.stack)-1] = diff
		case "*":
			mul := e.stack[len(e.stack)-1] * e.stack[len(e.stack)-2]
			e.stack = e.stack[:len(e.stack)-1]
			e.stack[len(e.stack)-1] = mul

		case "/":
			div := e.stack[len(e.stack)-2]
			if e.stack[len(e.stack)-1] == 0 {
				e.stack = tmp
				return errors.New("division by 0")
			}
			div /= e.stack[len(e.stack)-1]
			e.stack = e.stack[:len(e.stack)-1]
			e.stack[len(e.stack)-1] = div

		case "dup":
			e.stack = append(e.stack, e.stack[len(e.stack)-1])
		case "drop":
			e.stack = e.stack[:len(e.stack)-1]
		case "swap":
			e.stack[len(e.stack)-2], e.stack[len(e.stack)-1] =
				e.stack[len(e.stack)-1], e.stack[len(e.stack)-2]
		case "over":
			e.stack = append(e.stack, e.stack[len(e.stack)-2])
		default:
			num, _ := strconv.Atoi(str)
			e.stack = append(e.stack, num)
		}
	}

	return nil
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	args := strings.Fields(row)

	for i := range args {
		args[i] = strings.ToLower(args[i])
	}

	// this is definition
	if args[0] == ":" {
		tmp := make([]string, 0, len(args))

		if _, err := strconv.Atoi(args[1]); err == nil {
			return e.stack, errors.New("numbers redefinition")
		}

		for i := 2; i+1 < len(args); i++ {
			val, exist := e.m[args[i]]
			if exist {
				tmp = append(tmp, val...)
			} else {
				tmp = append(tmp, args[i])
			}
		}

		if args[len(args)-1] != ";" {
			return e.stack, errors.New("wrong syntax")
		}

		e.m[args[1]] = tmp
		return e.stack, nil
	}

	//otherwise this is a commands
	for _, str := range args {
		val, exist := e.m[str]

		if !exist {
			num, err := strconv.Atoi(str)
			if err != nil {
				return e.stack, errors.New("undefined variable")
			}
			e.stack = append(e.stack, num)
		} else {
			err := e.doCommand(val)
			if err != nil {
				return e.stack, err
			}
		}
	}
	return e.stack, nil
}
