package parser

import "testing"

func TestSymbols(t *testing.T) {
	table := new(symbols)
	table.openScope()
	table.add("a")
	table.add("b")
	table.add("c")
	table.add("d")
	table.openScope()
	table.add("e")
	table.add("f")
	table.openScope()
	table.add("a")
	table.add("b")

	b := table.find("d")
	if !b {
		t.Fail()
	}

	table.closeScope()
	table.closeScope()

	b = table.find("f")
	if b {
		t.Fail()
	}

	table.closeScope()

	if table.scopes != nil {
		t.Fail()
	}
}
