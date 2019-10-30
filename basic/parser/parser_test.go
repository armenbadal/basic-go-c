package parser

import "testing"

func oneCase(file string, t *testing.T) {
	pars := NewParser(file)
	if pars == nil {
		t.Fail()
		return
	}

	_, err := pars.Parse()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestT0(t *testing.T) {
	oneCase("../../examples/ex0.bas", t)
	oneCase("../../examples/ex1.bas", t)
	oneCase("../../examples/ex2.bas", t)
	oneCase("../../examples/ex3.bas", t)
	oneCase("../../examples/ex4.bas", t)
	oneCase("../../examples/ex5.bas", t)
	oneCase("../../examples/ex6.bas", t)
	oneCase("../../examples/ex7.bas", t)
}
