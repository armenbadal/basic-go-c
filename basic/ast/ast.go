package ast

// Node աբստրակտ քերականական ծառի հանգույց
type Node interface{}

// Number ...
type Number struct {
	Value int
}

// Text ...
type Text struct {
	Value string
}

// Variable ...
type Variable struct {
	Name  string
	Index int
}

// Unary ...
type Unary struct {
	Operation  string
	Expression Node
}

// Binary ...
type Binary struct {
	Operation   string
	Left, Right Node
}

// Apply ...
type Apply struct {
	Callee    *Subroutine
	Arguments []Node
}

// Let ...
type Let struct {
	VarName string
	Expr    Node
}

// Input ...
type Input struct {
	VarName string
}

// Print ...
type Print struct {
	Expr Node
}

// If ...
type If struct {
	Condition   Node
	Decision    Node
	Alternative Node
}

// While ...
type While struct {
	Condition Node
	Body      Node
}

// For ...
type For struct {
	Parameter string
	Begin     Node
	End       Node
	Step      Node
	Body      Node
}

// Call ...
type Call struct {
	Callee    *Subroutine
	Arguments []Node
}

// Sequence ...
type Sequence struct {
	Items []Node
}

// Subroutine ...
type Subroutine struct {
	Name       string
	Parameters []string
	Locals     []string
	Body       Node
}

// Program ...
type Program struct {
	Members map[string]*Subroutine
}
