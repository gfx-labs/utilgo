package lexer

type RuneStack struct {
	xs  []rune
	pos int
}

func (s *RuneStack) push(r rune) {
	s.xs = append(s.xs, r)
}

func (s *RuneStack) pop() (o rune) {
	if len(s.xs) == 0 {
		return EOFRune
	}
	s.xs, o = s.xs[:len(s.xs)-1], s.xs[len(s.xs)-1]
	return
}

func (s *RuneStack) clear() {
	s.xs = s.xs[:0]
}
