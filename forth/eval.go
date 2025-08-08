//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type IdentifierKind = byte

type IdentifierMapValue struct {
	kind     IdentifierKind
	value    []int
	function [](func(stack *[]int) error)
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
			kind:     function,
			value:    nil,
			function: [](func(stack *[]int) error){sumOperator},
		},
	}

	return &Evaluator{identifierMap}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	stack := []int{}

	words := strings.Split(row, " ")

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

// TODO: доделать функцию для определения идентфикатора
func (e *Evaluator) defineIdentifier(words []string, stack *[]int) error {
	var err error

	identifierName := words[0]

	l := len(words) - 1

	//Нужно проверить является ли первое слово функицей, если да то записываем массив функций в идентфикатор
	//Если нет, то 2 варианта: 1 - если не встречается функция, то записываем массив значений в идентфикатор
	//2 - если встретилась функция, то нужно вызвать computeStack и записать в идентифкатор верхнее значение

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
				funLen := len(identifierValue.function)

				for j := 0; j < funLen && err == nil; j++ {
					fun := identifierValue.function[j]

					err = fun(stack)
				}
			} else {
				for _, v := range identifierValue.value {
					*stack = append(*stack, v)
				}
			}
		} else if num, err := strconv.Atoi(word); err == nil {
			*stack = append(*stack, num)
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

func pop(stack *[]int) int {
	top := (*stack)[len(*stack)-1]

	*stack = (*stack)[:len(*stack)-1]

	return top
}
