package gun_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gfx-labs/utilgo/gun"
)

var exampleConfigOne struct {
	Field1    string `yaml:"field_one" env:"FIELD_ONE" json:"f_1" default:"IM DEFAULT HI"`
	Field2    int32
	ManyField []string
	Some      struct {
		Nested struct {
			Item string
		}
	}
}

func TestConfigOne(t *testing.T) {
	os.Setenv("SOME_NESTED_ITEM", "there")
	os.Setenv("FIELD_ONE", "hi")
	os.Setenv("FIELD_2", "69")
	os.Setenv("MANY_FIELD", "one,two,three")
	gun.LoadPrefix(&exampleConfigOne, "")

	fmt.Printf("t: %+v\n", exampleConfigOne)
}

func TestPrefixOne(t *testing.T) {
	os.Setenv("MANY_FIELD", "one,two,three")
	os.Setenv("TEST_MANY_FIELD", "four,five,six")
	os.Setenv("TEST2_MANY_FIELD", "seven,eight,nine")
	gun.LoadPrefix(&exampleConfigOne, "TEST")
	fmt.Printf("prefix_test: %+v\n", exampleConfigOne)
	gun.LoadPrefix(&exampleConfigOne, "TEST2")
	fmt.Printf("prefix_test2: %+v\n", exampleConfigOne)
}
