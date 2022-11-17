package main

import (
	"log"

	"gfx.cafe/util/go/lexer"
)

const (
	FloatToken lexer.TokenType = iota + 100
	ScientificToken
)

func LexExpr(l *lexer.Lex) lexer.StateFn {
	for {
		switch l.Cur() {
		case "":
		case `'`, `"`:
			return LexEscapableString(rune(l.Cur()[0]))
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			return LexNumeric
		case " ":
			l.Ignore()
		default:
			l.Emit(lexer.SymbolToken)
		}
		if l.Next() == lexer.EOFRune {
			break
		}
	}
	l.Emit(lexer.EOFToken)
	return nil
}

func LexNumeric(l *lexer.Lex) lexer.StateFn {
	for {
		switch l.Peek() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return LexNumeric
		case 'e':
			return LexScientific
		case '.':
			return LexFloat
		default:
			l.Emit(lexer.IntegerToken)
			return LexExpr
		}
	}
}

func LexFloat(l *lexer.Lex) lexer.StateFn {
	l.Next()
	for {
		switch l.Peek() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.TakeWhile("0123456789")
			l.Emit(FloatToken)
			return nil
		default:
			l.Error("bad numeric" + l.Cur())
		}
		if l.Next() == lexer.EOFRune {
			break
		}
	}
	return nil
}
func LexScientific(l *lexer.Lex) lexer.StateFn {
	l.Next()
	for {
		switch l.Peek() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.TakeWhile("0123456789")
			l.Emit(ScientificToken)
			return LexExpr
		default:
			l.Error("bad numeric" + l.Cur())
		}
		if l.Next() == lexer.EOFRune {
			break
		}
	}
	return nil
}

func LexEscapableString(end rune) func(*lexer.Lex) lexer.StateFn {
	return func(l *lexer.Lex) lexer.StateFn {
		l.Ignore()
		escape := false
		for {
			switch l.Peek() {
			case '\\':
				escape = true
			default:
				if !escape {
					switch l.Peek() {
					case end:
						l.Emit(lexer.StringToken)
						l.Next()
						l.Ignore()
						return LexExpr
					}
				}
				l.Next()
				escape = false
			}
			if l.Next() == lexer.EOFRune {
				l.Emit(lexer.EOFToken)
				l.Error("unexpected EOF")
				return nil
			}
		}
	}
}

func main() {
	expr := `1e44 + "swa \" gg"`
	lexer.New(expr, LexExpr).ConsumeWith(func(tok *lexer.Token) {
		log.Println(tok.Value)
	})
	_ = expr
}
