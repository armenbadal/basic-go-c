package parser

import (
	"basic/ast"
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
)

// վերլուծված ենթածրագրերի «գլոբալ» ցուցակ
var subrs = map[string]*ast.Subroutine{}

// անհայտ ենթածրագրերի կանչերի ցուցակ
var clinks = map[string]*list.List{}

// Parser Շարահյուսական վերլուծիչի ստրուկտուրան։
type Parser struct {
	scer      *scanner
	lookahead *lexeme

	program *ast.Program
}

// NewParser Ստեղծում և վերադարձնում է շարահյուսական վերլուծիչի նոր օբյեկտ։
func NewParser(filename string) *Parser {
	// բացել ֆայլային հոսքը
	rd, er := os.Open(filename)
	if er != nil {
		// TODO: signal for failure
		println("Cannot open file")
		return nil
	}
	defer rd.Close()

	// ստեղծել շարահյուսական վերլուծիչի օբյեկտը
	pars := new(Parser)
	pars.scer = new(scanner)
	pars.scer.source = bufio.NewReader(rd)
	pars.scer.line = 1
	pars.scer.read()
	pars.lookahead = pars.scer.next()

	pars.program = &ast.Program{Members: make(map[string]*ast.Subroutine)}

	return pars
}

// Parse Վերլուծությունը սկսող արտաքին ֆունկցիա
func (p *Parser) Parse() (*ast.Program, error) {
	p.parseProgram()
	return p.program, nil
}

// Վերլուծել ամբողջ ծրագիրը.
//
// Program = { Subroutine }.
//
func (p *Parser) parseProgram() error {
	// բաց թողնել ֆայլի սկզբի դատարկ տողերը
	for p.has(xNewLine) {
		p.lookahead = p.scer.next()
	}

	for p.has(xSubroutine) {
		p0, err := p.parseSubroutine()
		if err != nil {
			return err
		}

		p.program.Members[p0.Name] = p0

		if err = p.parseNewLines(); err != nil {
			return err
		}
	}

	return nil
}

// Վերլուծել նոր տողերի նիշերի հաջորդականությունը
//
// NewLines = NEWLINE { NEWLINE }.
//
func (p *Parser) parseNewLines() error {
	if _, err := p.match(xNewLine); err != nil {
		return err
	}

	for p.lookahead.is(xNewLine) {
		if _, err := p.match(xNewLine); err != nil {
			return err
		}
	}

	return nil
}

// Վերլուծել ենթածրագիրը
//
// Subroutine = 'SUB' IDENT ['(' [IDENT {',' IDENT}] ')'] NewLines
//              { Statement NewLines } 'END' SUB'.
//
func (p *Parser) parseSubroutine() (*ast.Subroutine, error) {
	p.match(xSubroutine)
	name, err := p.match(xIdent)
	if err != nil {
		return nil, err
	}
	pars := make([]string, 0, 16)
	if p.has(xLeftPar) {
		p.match(xLeftPar)
		if p.has(xIdent) {
			pnm, _ := p.match(xIdent)
			pars = append(pars, pnm)
			for p.has(xComma) {
				p.match(xComma)
				pnm, err := p.match(xIdent)
				if err != nil {
					return nil, err
				}
				pars = append(pars, pnm)
			}
		}
		if _, err := p.match(xRightPar); err != nil {
			return nil, err
		}
	}

	sub := &ast.Subroutine{
		Name:       name,
		Parameters: pars,
		Locals:     nil,
		Body:       nil}
	// TODO: add this subroutine to the program
	// TODO: add parapeters and subroutine name to the locals

	body, err := p.parseSequence()
	if err != nil {
		return nil, err
	}
	sub.Body = body

	if _, err := p.match(xEnd); err != nil {
		return nil, err
	}
	if _, err := p.match(xSubroutine); err != nil {
		return nil, err
	}

	// subrs[name] = sub

	// if clinks[name] != nil {
	// 	for e := clinks[name].Front(); e != nil; e = e.Next() {
	// 		switch coa := e.Value.(type) {
	// 		case ast.Call:
	// 			coa.SetCallee(sub)
	// 		case ast.Apply:
	// 			coa.SetCallee(sub)
	// 		}
	// 	}
	// 	delete(clinks, name)
	// }

	return sub, nil
}

// Վերլուծել հրամանների հաջորդականություն
func (p *Parser) parseSequence() (ast.Node, error) {
	if err := p.parseNewLines(); err != nil {
		return nil, err
	}

	res := make([]ast.Node, 0, 16)

	for p.has(xLet, xInput, xPrint, xIf, xWhile, xFor, xCall) {
		var stat ast.Node
		switch {
		case p.lookahead.is(xLet):
			stat, _ = p.parseLet()
		case p.lookahead.is(xInput):
			stat, _ = p.parseInput()
		case p.lookahead.is(xPrint):
			stat, _ = p.parsePrint()
		case p.lookahead.is(xIf):
			stat, _ = p.parseIf()
		case p.lookahead.is(xWhile):
			stat, _ = p.parseWhile()
		case p.lookahead.is(xFor):
			stat, _ = p.parseFor()
		case p.lookahead.is(xCall):
			stat, _ = p.parseCall()
		}

		if err := p.parseNewLines(); err != nil {
			return nil, err
		}

		res = append(res, stat)
	}

	return res, nil
}

// Վերլուծել վերագրման հրամանը
//
// Statement = 'LET' IDENT '=' Expression.
//
func (p *Parser) parseLet() (ast.Node, error) {
	p.match(xLet)

	vn, err := p.match(xIdent)
	if err != nil {
		return nil, err
	}

	if _, err = p.match(xEq); err != nil {
		return nil, err
	}

	e0, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	return &ast.Let{VarName: vn, Expr: e0}, nil
}

// Ներմուծման հրամանի վերլուծությունը.
//
// Statement = 'INPUT' IDENT.
//
func (p *Parser) parseInput() (ast.Node, error) {
	p.match(xInput)
	nam, err := p.match(xIdent)
	if err != nil {
		return nil, err
	}
	return &ast.Input{VarName: nam}, nil
}

// Արտածման հրամանի վերլուծությունը.
//
// Statement = 'PRINT' Expression.
//
func (p *Parser) parsePrint() (ast.Node, error) {
	p.match(xPrint)
	e0, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	return &ast.Print{Expr: e0}, nil
}

// Ճյուղավորման հրամանի վերլուծությունը.
//
// Statement = 'IF' Expression 'THEN' NewLines { Statement NewLines }
//             { 'ELSEIF' Expression 'THEN' NewLines { Statement NewLines } }
//             [ 'ELSE' NewLines { Statement NewLines } ]
//             'END' 'IF'.
//
func (p *Parser) parseIf() (ast.Node, error) {
	p.match(xIf)

	c0, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.match(xThen); err != nil {
		return nil, err
	}

	s0, err := p.parseSequence()
	if err != nil {
		return nil, err
	}

	res := &ast.If{Condition: c0, Decision: s0, Alternative: nil}
	ipe := res
	for p.has(xElseIf) {
		p.match(xElseIf)
		c1, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.match(xThen); err != nil {
			return nil, err
		}
		s1, err := p.parseSequence()
		if err != nil {
			return nil, err
		}
		alt := &ast.If{Condition: c1, Decision: s1, Alternative: nil}
		ipe.Alternative = alt
		ipe = ipe.Alternative.(*ast.If)
	}
	if p.has(xElse) {
		p.match(xElse)
		s2, _ := p.parseSequence()
		ipe.Alternative = s2
	}

	if _, err := p.match(xEnd); err != nil {
		return nil, err
	}
	if _, err := p.match(xIf); err != nil {
		return nil, err
	}

	return res, nil
}

// Նախապայմանով ցիկլի վերլուծությունը
//
// Statement = 'WHILE' Expression NewLines
//             { Statement NewLines } 'END' 'WHILE'.
//
func (p *Parser) parseWhile() (ast.Node, error) {
	p.match(xWhile)
	c0, _ := p.parseExpression()
	b0, _ := p.parseSequence()
	p.match(xEnd)
	p.match(xWhile)
	return &ast.While{Condition: c0, Body: b0}, nil
}

// Պարամետրով ցիկլի վերլուծությունը
//
// Statement = 'FOR' IDENT '=' Expression 'TO' Expression
//             ['STEP' ['+'|'-'] NUMBER { Statement NewLines }
//             'END' 'FOR'.
//
func (p *Parser) parseFor() (ast.Node, error) {
	p.match(xFor)
	param := p.lookahead.value
	p.match(xIdent)
	p.match(xEq)
	b0, _ := p.parseExpression()
	p.match(xTo)
	e0, _ := p.parseExpression()
	var num int64 = 1
	if p.has(xStep) {
		p.match(xStep)
		posi := true
		if p.has(xSub) {
			p.match(xSub)
			posi = false
		} else if p.has(xAdd) {
			p.match(xAdd)
		}
		lex := p.lookahead.value
		p.match(xNumber)
		num, _ = strconv.ParseInt(lex, 10, 32)
		if !posi {
			num = -num
		}
	}
	s0 := &ast.Number{Value: int(num)}
	dy, _ := p.parseSequence()

	if _, err := p.match(xEnd); err != nil {
		return nil, err
	}
	if _, err := p.match(xFor); err != nil {
		return nil, err
	}

	return &ast.For{
		Parameter: param,
		Begin:     b0,
		End:       e0,
		Step:      s0,
		Body:      dy}, nil
}

// Ենթածրագրի կանչի վերլուծությունը
//
// Statement = 'CALL' IDENT [Expression {',' Expression}].
//
func (p *Parser) parseCall() (ast.Node, error) {
	p.match(xCall)
	name, err := p.match(xIdent)
	if err != nil {
		return nil, err
	}

	args := make([]ast.Node, 0, 16)
	if p.has(xNumber, xText, xIdent, xSub, xNot, xLeftPar) {
		e0, _ := p.parseExpression()
		args = append(args, e0)
		for p.has(xComma) {
			p.match(xComma)
			e1, _ := p.parseExpression()
			args = append(args, e1)
		}
	}

	sp, defined := subrs[name]
	if defined {
		return &ast.Call{Callee: sp, Arguments: args}, nil
	}

	dummy := &ast.Subroutine{Name: "__dummy__", Parameters: nil, Body: nil}
	dcall := &ast.Call{Callee: dummy, Arguments: args}
	if clinks[name] == nil {
		clinks[name] = list.New()
	}
	clinks[name].PushBack(dcall)

	return dcall, nil
}

//
func (p *Parser) parseExpression() (ast.Node, error) {
	res, _ := p.parseConjunction()
	for p.has(xOr) {
		p.match(xOr)
		e0, _ := p.parseConjunction()
		res = &ast.Binary{Operation: "OR", Left: res, Right: e0}
	}
	return res, nil
}

//
func (p *Parser) parseConjunction() (ast.Node, error) {
	res, _ := p.parseEquality()
	for p.has(xAnd) {
		p.match(xAnd)
		e0, _ := p.parseEquality()
		res = &ast.Binary{Operation: "AND", Left: res, Right: e0}
	}
	return res, nil
}

//
func (p *Parser) parseEquality() (ast.Node, error) {
	res, _ := p.parseComparison()
	if p.has(xEq, xNe) {
		var opc string
		switch p.lookahead.token {
		case xEq:
			opc = "EQ"
			p.match(xEq)
		case xNe:
			opc = "NE"
			p.match(xNe)
		}
		e0, _ := p.parseComparison()
		res = &ast.Binary{Operation: opc, Left: res, Right: e0}
	}

	return res, nil
}

//
func (p *Parser) parseComparison() (ast.Node, error) {
	res, _ := p.parseAddition()
	if p.has(xGt, xGe, xLt, xLe) {
		var opc string
		switch p.lookahead.token {
		case xGt:
			opc = "GT"
			p.match(xGt)
		case xGe:
			opc = "GE"
			p.match(xGe)
		case xLt:
			opc = "LT"
			p.match(xLt)
		case xLe:
			opc = "LE"
			p.match(xLe)
		}
		e0, _ := p.parseAddition()
		res = &ast.Binary{Operation: opc, Left: res, Right: e0}
	}
	return res, nil
}

//
func (p *Parser) parseAddition() (ast.Node, error) {
	res, _ := p.parseMultiplication()
	for p.has(xAdd, xSub, xAmp) {
		var opc string
		switch p.lookahead.token {
		case xAdd:
			opc = "ADD"
			p.match(xAdd)
		case xSub:
			opc = "SUB"
			p.match(xSub)
		case xAmp:
			opc = "CONC"
			p.match(xAmp)
		}
		e0, _ := p.parseMultiplication()
		res = &ast.Binary{Operation: opc, Left: res, Right: e0}
	}

	return res, nil
}

//
func (p *Parser) parseMultiplication() (ast.Node, error) {
	res, _ := p.parsePower()
	for p.has(xMul, xDiv, xMod) {
		opc := ""
		switch p.lookahead.token {
		case xMul:
			opc = "MUL"
			p.match(xMul)
		case xDiv:
			opc = "DIV"
			p.match(xDiv)
		case xMod:
			opc = "MOD"
			p.match(xMod)
		}
		e0, _ := p.parsePower()
		res = &ast.Binary{Operation: opc, Left: res, Right: e0}
	}

	return res, nil
}

// Ատիճան բարձրացնելու գործողությունը
//
// Power = Factor '^' Power.
//
func (p *Parser) parsePower() (ast.Node, error) {
	res, _ := p.parseFactor()
	if p.has(xPow) {
		p.match(xPow)
		e0, _ := p.parsePower()
		res = &ast.Binary{Operation: "POW", Left: res, Right: e0}
	}

	return res, nil
}

// Պարզագույն արտահայտությունների վերլուծությունը
//
// Factor = NUMBER | TEXT | IDENT
//        | SUB Factor
//        | NOT Factor
//        | '(' Expression ')'.
//
func (p *Parser) parseFactor() (ast.Node, error) {
	switch p.lookahead.token {
	case xNumber:
		return p.parseNumber()
	case xText:
		return p.parseText()
	case xIdent:
		return p.parserVariableOrApply()
	case xSub, xNot:
		return p.parseUnary()
	case xLeftPar:
		return p.parseParenthesis()
	}

	return nil, nil
}

// թվային լիտերալ (ամբողջաթիվ)
func (p *Parser) parseNumber() (*ast.Number, error) {
	lex, _ := p.match(xNumber)
	val, _ := strconv.ParseInt(lex, 10, 32)
	return &ast.Number{Value: int(val)}, nil
}

// տեքստային լիտերալ
func (p *Parser) parseText() (*ast.Text, error) {
	val, _ := p.match(xText)
	return &ast.Text{Value: val}, nil
}

// փոփոխական կամ ֆունկցիայի կիրառություն
func (p *Parser) parserVariableOrApply() (ast.Node, error) {
	name, _ := p.match(xIdent)
	if p.has(xLeftPar) {
		return p.parseApply()
	}

	return &ast.Variable{Name: name}, nil
}

// ֆունկցիայի կանչ
func (p *Parser) parseApply() (ast.Node, error) {
	p.match(xLeftPar)
	args := make([]ast.Node, 0, 16)
	if p.has(xNumber, xText, xIdent, xSub, xNot, xLeftPar) {
		e0, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, e0)
		for p.has(xComma) {
			p.match(xComma)
			e1, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, e1)
		}
	}
	if _, err := p.match(xRightPar); err != nil {
		return nil, err
	}

	return &ast.Apply{Callee: nil, Arguments: args}, nil
}

// ունար գործողություն
func (p *Parser) parseUnary() (ast.Node, error) {
	var opc string
	if p.lookahead.is(xSub) {
		opc = "NEG"
		p.match(xSub)
	} else if p.lookahead.is(xNot) {
		opc = "NOT"
		p.match(xNot)
	}
	res, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	res = &ast.Unary{Operation: opc, Expression: res}
	return res, nil
}

// փակագծեր՝ բարձր առաջնահերթությամբ արտահայտություն
func (p *Parser) parseParenthesis() (ast.Node, error) {
	p.match(xLeftPar)
	res, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if _, err = p.match(xRightPar); err != nil {
		return nil, err
	}
	return res, nil
}

//
func (p *Parser) has(tokens ...int) bool {
	return p.lookahead.is(tokens...)
}

//
func (p *Parser) match(exp int) (string, error) {
	if p.lookahead.is(exp) {
		lex := p.lookahead.value
		p.lookahead = p.scer.next()
		return lex, nil
	}

	return "", fmt.Errorf("Տող %d: Վերլուծության սխալ", p.lookahead.line)
}
