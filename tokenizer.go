package simript

import (
	"regexp"
)

type Token struct {
	token_type string
	raw				 string
	newLineable bool
}


type patternF func([]string) Token

func GetBracketDirection(bracket rune) bool {
	return (bracket == '(') || (bracket == '[') || (bracket == '{')
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}else{ return 0 }
}

func strContentPattern(c rune) string {
	return `(?:[^\r\n\\`+string(c)+`]|\\t|\\r|\\n|\\x\h{1,2}|\\u\h{1,4}|\\u\{\h+\}|\\.)`
}

var TokenPatterns map[*regexp.Regexp]patternF = map[*regexp.Regexp]patternF {
	// ケツカンマしないと怒られるので注意
	regexp.MustCompile(`^[$\p{L}_]+[$\p{L}\d_]*`):
		func (matches []string) Token {
			return Token{"IDENTIFIER",matches[0],false}
		},
	regexp.MustCompile(`^[-+]?(?:[1-9][0-9_]+|[0-9])(\.[0-9](?:[0-9]+|[0-9_]+[0-9])?)(?:e[-+]?\d+)?`):
		func (matches []string) Token {
			return Token{"NUMBER",matches[0],false}
		},
	regexp.MustCompile(`^"`+strContentPattern('"')+`"|'`+strContentPattern('\'')+`'`):
		func (matches []string) Token {
			return Token{"STRING",matches[0],false}
		},
	regexp.MustCompile(`^(?:(?:!=|==|!|=)|(?:&&|\|\||\*\*|[+\-*/^~&|])=?|\.)`):
		func (matches []string) Token {
			return Token{"OPERATOR",matches[0],true}
		},
	regexp.MustCompile(`^[()]`):
		func (matches []string) Token {
			return Token{"BRACKET_" + ([]string{"LEFT","RIGHT"}[BoolToInt(GetBracketDirection([]rune(matches[0])[0]))]), matches[0],true}
		},
}

func Tokenize (code string) []Token {
	tokens := []Token{}
	var prevNewLineable bool = true
	for i := 0; i<len(code); {
		for k,v := range TokenPatterns {
			m := k.FindStringSubmatch(code[i:])
			if m == nil && !prevNewLineable && code[i:i+1] == "\n" {
				return nil
			}
			if m == nil {
				return nil
			}
			i += len(m[0])
			tokens = append(tokens,v(m))
		}
	}
	return tokens
}