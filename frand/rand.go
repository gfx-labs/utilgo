package frand

const (
	// Set of characters to use for generating random strings
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
	Ascii        = Alphanumeric + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

// IntRange returns a random integer in the range from min to max.
func IntRange(min, max int) int {
	var result int
	switch {
	case max == min:
		result = max
	case min > max:
		min, max = max, min
		fallthrough
	case max > min:
		maxRand := max - min
		b := Intn(maxRand)
		result = min + b
	}
	return result
}

// String returns a random string n characters long, composed of entities
// from charset.
func String(n int, charset string) (string, error) {
	randstr := make([]byte, n) // Random string to return
	for i := 0; i < n; i++ {
		b := Intn(len(charset))
		r := int(b)
		randstr[i] = charset[r]
	}
	return string(randstr), nil
}

// StringRange returns a random string at least min and no more than max
// characters long, composed of entitites from charset.
func StringRange(min, max int, charset string) (string, error) {
	//
	// First determine the length of string to be generated
	//
	var err error      // Holds errors
	var strlen int     // Length of random string to generate
	var randstr string // Random string to return
	strlen = IntRange(min, max)
	randstr, err = String(strlen, charset)
	if err != nil {
		return randstr, err
	}
	return randstr, nil
}

// AlphaRange returns a random alphanumeric string at least min and no more
// than max characters long.
func AlphaStringRange(min, max int) (string, error) {
	return StringRange(min, max, Alphanumeric)
}

// AlphaString returns a random alphanumeric string n characters long.
func AlphaString(n int) (string, error) {
	return String(n, Alphanumeric)
}

// ChoiceAny returns a random selection from an array of T.
func ChoiceAny[T any](choices []T) T {
	length := len(choices)
	i := IntRange(0, length)
	return choices[i]
}

// ChoiceString returns a random selection from an array of strings.
func ChoiceString(choices []string) string {
	length := len(choices)
	i := IntRange(0, length)
	return choices[i]
}

// ChoiceInt returns a random selection from an array of integers.
func ChoiceInt(choices []int) int {
	length := len(choices)
	i := IntRange(0, length)
	return choices[i]
}

type Choice = ChoiceT[any]

func WeightedChoice(choices []Choice) Choice {
	return WeightedChoiceT(choices)
}

// A Choice contains a generic item and a weight controlling the frequency with
// which it will be selected.
type ChoiceT[T any] struct {
	Weight int
	Item   T
}

func WeightedChoiceT[T any](choices []ChoiceT[T]) ChoiceT[T] {
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r := IntRange(0, sum)
	var c ChoiceT[T]
	for _, c = range choices {
		r -= c.Weight
		if r < 0 {
			return c
		}
	}
	panic("should not reach here. please report bug in gfx.cafe/util/go")
}

func InvertChoicesT[T any](choices []ChoiceT[T]) []ChoiceT[T] {
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	sum2 := 0
	for idx, v := range choices {
		tmp := 100 * int(float64(sum)/(float64(v.Weight)+0.001))
		sum2 = sum2 + tmp
		choices[idx].Weight = tmp
	}
	return choices
}
