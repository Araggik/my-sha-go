//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type IdentifierKind = byte

type IdentifierMapValue struct {
	kind      IdentifierKind
	values    []int
	functions [](func(stack *[]int) error)
}

type Evaluator struct {
	identifierMap map[string]IdentifierMapValue
}

const (
	value IdentifierKind = iota
	function
)

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	identifierMap := map[string]IdentifierMapValue{
		"+": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){sumOperator},
		},
		"-": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){minusOperator},
		},
		"*": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){factorOperator},
		},
		"/": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){deriveOperator},
		},
		"dup": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){dupFunc},
		},
		"over": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){overFunc},
		},
		"drop": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){dropFunc},
		},
		"swap": {
			kind:      function,
			values:    nil,
			functions: [](func(stack *[]int) error){swapFunc},
		},
	}

	return &Evaluator{identifierMap}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	stack := []int{}

	lowerRow := strings.ToLower(row)

	words := strings.Split(lowerRow, " ")

	var err error

	if words[0] == ":" {
		lastIndex := len(words) - 1

		if words[lastIndex] == ";" {
			err = e.defineIdentifier(words[1:lastIndex], &stack)
		} else {
			err = errors.New("в конце нет ;")
		}
	} else {
		err = e.computeStack(words, &stack)
	}

	return stack, err
}

func (e *Evaluator) defineIdentifier(words []string, stack *[]int) error {
	var err error

	identifierName := words[0]

	if _, nameConvErr := strconv.Atoi(identifierName); nameConvErr == nil {
		err = errors.New("имя идентификатора не может быть числом")
	}

	l := len(words)

	//Нужно проверить является ли первое слово функицей, если да то записываем массив функций в идентфикатор
	//Если нет, то 2 варианта: 1 - если не встречается функция, то записываем массив значений в идентфикатор
	//2 - если встретилась функция, то нужно вызвать computeStack и записать в идентифкатор верхнее значение
	if val, ok := e.identifierMap[words[1]]; ok && val.kind == function {
		var funcArr [](func(stack *[]int) error)

		for i := 1; i < l && err == nil; i++ {
			if identifierValue, ok := e.identifierMap[words[i]]; ok && identifierValue.kind == function {
				funcArr = append(funcArr, identifierValue.functions...)
			} else {
				err = errors.New("неправильное определение идентификатора для функции")
			}
		}

		if err == nil {
			e.identifierMap[identifierName] = IdentifierMapValue{function, nil, funcArr}
		}
	} else {
		isFindFun := false

		for i := 1; i < l && err == nil && !isFindFun; i++ {
			word := words[i]

			if identifierValue, ok := e.identifierMap[word]; ok {
				if identifierValue.kind == function {
					isFindFun = true

					err = e.computeStack(words[i:], stack)
				} else {
					*stack = append(*stack, identifierValue.values...)
				}
			} else {
				if num, convErr := strconv.Atoi(word); convErr == nil {
					*stack = append(*stack, num)
				} else {
					err = convErr
				}
			}
		}

		if err == nil {
			var values []int

			if isFindFun {
				values = []int{(*stack)[len(*stack)-1]}
			} else {
				values = *stack
			}

			e.identifierMap[identifierName] = IdentifierMapValue{value, values, nil}
		}
	}

	return err
}

func (e *Evaluator) computeStack(words []string, stack *[]int) error {
	var err error

	l := len(words)

	for i := 0; i < l && err == nil; i++ {
		word := words[i]

		identifierValue, ok := e.identifierMap[word]

		if ok {
			kind := identifierValue.kind

			if kind == function {
				funLen := len(identifierValue.functions)

				for j := 0; j < funLen && err == nil; j++ {
					fun := identifierValue.functions[j]

					err = fun(stack)
				}
			} else {
				*stack = append(*stack, identifierValue.values...)
			}
		} else {
			if num, convErr := strconv.Atoi(word); convErr == nil {
				*stack = append(*stack, num)
			} else {
				err = convErr
			}
		}
	}

	return err
}

func sumOperator(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для сложения")
	}

	last := pop(stack)
	prev := pop(stack)

	*stack = append(*stack, last+prev)

	return nil
}

func minusOperator(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для вычитания")
	}

	last := pop(stack)
	prev := pop(stack)

	*stack = append(*stack, prev-last)

	return nil
}

func factorOperator(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для умножения")
	}

	last := pop(stack)
	prev := pop(stack)

	*stack = append(*stack, prev*last)

	return nil
}

func deriveOperator(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для деления")
	}

	last := pop(stack)

	if last == 0 {
		return errors.New("деление на ноль")
	}

	prev := pop(stack)

	*stack = append(*stack, prev/last)

	return nil
}

func dupFunc(stack *[]int) error {
	if len(*stack) < 1 {
		return errors.New("не хватает аргументов для dup")
	}

	last := (*stack)[len(*stack)-1]

	*stack = append(*stack, last)

	return nil
}

func overFunc(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для over")
	}

	prev := (*stack)[len(*stack)-2]

	*stack = append(*stack, prev)

	return nil
}

func dropFunc(stack *[]int) error {
	if len(*stack) < 1 {
		return errors.New("не хватает аргументов для drop")
	}

	pop(stack)

	return nil
}

func swapFunc(stack *[]int) error {
	if len(*stack) < 2 {
		return errors.New("не хватает аргументов для swap")
	}

	l := len(*stack)

	last := (*stack)[l-1]

	(*stack)[l-1] = (*stack)[l-2]
	(*stack)[l-2] = last

	return nil
}

func pop(stack *[]int) int {
	top := (*stack)[len(*stack)-1]

	*stack = (*stack)[:len(*stack)-1]

	return top
}
