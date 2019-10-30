package ast

// Node աբստրակտ քերականական ծառի հանգույց
type Node interface{}

type Number struct {
	Value int
}

type Text struct {
	Value string
}

type Variable struct {
	Name string
}

type Unary struct {
	Operation  string
	Expression Node
}

type Binary struct {
	Operation   string
	Left, Right Node
}

type Apply struct {
	Callee    *Subroutine
	Arguments []Node
}

type Let struct {
	VarName string
	Expr    Node
}

type Input struct {
	VarName string
}

type Print struct {
	Expr Node
}

type If struct {
	Condition   Node
	Decision    Node
	Alternative Node
}

type While struct {
	Condition Node
	Body      Node
}

type For struct {
	Parameter string
	Begin     Node
	End       Node
	Step      Node
	Body      Node
}

type Call struct {
	Callee    *Subroutine
	Arguments []Node
}

type Sequence struct {
	Items []Node
}

type Subroutine struct {
	Name       string
	Parameters []string
	Body       Node
}

type Program struct {
	Members map[string]*Subroutine
}
