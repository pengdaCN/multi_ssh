package lexer

type (
	TokenKind int
	Token     struct {
		filename string
		line     string
		kind     TokenKind
		// 值的字符串表示
		val string
		// 值的整数表示
		valI int
		// 值的浮点数表示
		valF float64
	}
)

const (
	TOKEN_EOF        TokenKind = iota // end-of-file
	TOKEN_SEP_LPAREN                  // (
	TOKEN_SEP_RPAREN                  // )
	TOKEN_OP_MINUS                    // -
	TOKEN_OP_ASSIGN                   // =
	TOKEN_OP_LE                       // <=
	TOKEN_OP_GT                       // >
	TOKEN_OP_GE                       // >=
	TOKEN_OP_EQ                       // ==
	TOKEN_OP_NE                       // !=
	TOKEN_OP_AND                      // and
	TOKEN_OP_OR                       // or
	TOKEN_OP_NOT                      // not
)

var (
	keywords = map[string]TokenKind{
		"and": TOKEN_OP_AND,
		"or":  TOKEN_OP_OR,
		"not": TOKEN_OP_NOT,
	}
)
