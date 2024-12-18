package main

import (
	"fmt"
	"slices"
)

type Program struct {
	a              int
	b              int
	c              int
	ops            []int
	programPointer int
	output         []int
}
type queueVal struct {
	offset int
	val    int
}

func main() {
	p := read()
	p.execute()
	fmt.Println(p.output)
	fmt.Println(solve(p))
}

func solve(p Program) int {
	queue := []queueVal{{len(p.ops) - 1, 0}}
	for len(queue) > 0 {
		el := queue[0]
		queue = append([]queueVal{}, queue[1:]...)
		for i := range 8 {
			newVal := el.val<<3 + i
			newProgram := Program{
				newVal,
				0,
				0,
				p.ops,
				0,
				[]int{},
			}
			newProgram.execute()
			if slices.Compare(p.ops[el.offset:], newProgram.output[:]) == 0 {
				if el.offset == 0 {
					return newVal
				}
				queue = append(queue, queueVal{el.offset - 1, newVal})
			}
		}
	}
	panic("no value")
}

func (p *Program) execute() {
	for p.programPointer < len(p.ops)-1 {
		op, operand := p.ops[p.programPointer], p.ops[p.programPointer+1]
		switch op {
		case 0:
			p.adv(p.combo_op(operand))
			p.programPointer += 2
		case 1:
			p.bxl(operand)
			p.programPointer += 2
		case 2:
			p.bst(p.combo_op(operand))
			p.programPointer += 2
		case 3:
			if !p.jnz(operand) {
				p.programPointer += 2
			}
		case 4:
			p.bxc(operand)
			p.programPointer += 2
		case 5:
			val := p.out(p.combo_op(operand))
			p.programPointer += 2
			p.output = append(p.output, val)
		case 6:
			p.bdv(p.combo_op(operand))
			p.programPointer += 2
		case 7:
			p.cdv(p.combo_op(operand))
			p.programPointer += 2
		}
	}
}

func (p *Program) combo_op(operand int) int {
	switch operand {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 3
	case 4:
		return p.a
	case 5:
		return p.b
	case 6:
		return p.c
	}
	panic("Unsopported operator")
}

func (p *Program) adv(operand int) {
	p.a = p.a / (1 << operand)
}

func (p *Program) bxl(operand int) {
	p.b = p.b ^ operand
}

func (p *Program) bst(operand int) {
	p.b = operand % 8
}

func (p *Program) jnz(operand int) bool {
	if p.a == 0 {
		return false
	}
	p.programPointer = operand
	return true
}

func (p *Program) bxc(operand int) {
	p.b = p.b ^ p.c
}

func (p *Program) out(operand int) int {
	return operand % 8
}

func (p *Program) bdv(operand int) {
	p.b = p.a / (1 << operand)
}

func (p *Program) cdv(operand int) {
	p.c = p.a / (1 << operand)
}

func read() Program {
	p := Program{
		17323786,
		0,
		0,
		[]int{2, 4, 1, 1, 7, 5, 1, 5, 4, 1, 5, 5, 0, 3, 3, 0},
		0,
		[]int{},
	}
	return p
}
