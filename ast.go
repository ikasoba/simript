package simript

import (
	"errors"
	"strconv"
)

type any = interface{}

type ASTTypes int
type ValueType int

const (
	INVOKE_FUNC ASTTypes = iota
	AST_ROOT
	PRIMITIVE_VALUE
)
const (
	VALUE_NUMBER ValueType = iota
	VALUE_STRING
)

type Value struct {
	valueType ValueType
	value     any
}

type ASTInterface interface {
	getType() ASTTypes
}

type ASTInvokeFunction struct {
	args []Value
	funcName string
}

type ASTValue struct {
	value *Value
}

type ASTRoot struct {
	children []ASTInterface
}

func (root ASTRoot) getType() ASTTypes {
	return AST_ROOT
}

func (ast ASTInvokeFunction) getType() ASTTypes {
	return INVOKE_FUNC
}

/*
	関数名を第一引数、そのあとの`(`以降のトークンを第二引数へ
*/
func parseInvokeFunctions(funcName string, tokens []Token) (ASTInterface, error) {
	var args []Value
	for i:=0; i<len(tokens); i++ {
		t := tokens[i]
		if t.token_type == "NUMBER" || t.token_type == "STRING" {
			var v Value
			if t.token_type == "NUMBER" {
				n,e := strconv.ParseFloat(t.raw, 64)
				if e != nil {
					return nil, e
				}
				v = Value{VALUE_NUMBER, n}
			}
			if t.token_type == "STRING" {
				v = Value{VALUE_NUMBER, t.raw}
			}
			args = append(args,v)
		}
	}
	return ASTInvokeFunction{args, funcName}, nil
}

func CreateASTFromTokens(tokens []Token) (ASTRoot,error) {
	var root ASTRoot = ASTRoot{[]ASTInterface{}}
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if t.token_type == "IDENTIFIER" && (i+1<len(tokens) && tokens[i+1].token_type == "BRACKET_LEFT") {
			i += 2
			if i >= len(tokens) {
				return root, errors.New("invalid or unexpected token")
			}
			invokeTree, err := parseInvokeFunctions(t.raw,tokens[i:])
			if err != nil {
				return root, err
			}
			root.children = append(root.children,invokeTree)
		}
	}
	return root, nil
}
