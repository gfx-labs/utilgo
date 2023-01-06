package gun_test

import (
	"fmt"
	"os"
	"testing"

	"gfx.cafe/util/go/gun"
)

var exampleConfigOne struct {
	Field1     string `yaml:"field_one"`
	Field2     int32
	MANY_FIELD []string
}

func TestConfigOne(t *testing.T) {
	os.Setenv("MANY_FIELD", "one,two,three")
	gun.Load(&exampleConfigOne)
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
