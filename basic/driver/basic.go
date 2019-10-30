package main

import (
	"basic/ast"
)

func main() {
	var a0 ast.Number
	a0.Value = 777
	var a1 ast.Text
	a1.Value = "Halo"

	println("BASIC-S Compiler in Go")

	println(a0.Value)
	println(a1.Value)

	//pars := parser.NewParser("../examples/ex0.bas")
}
