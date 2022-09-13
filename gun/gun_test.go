package gun_test

import (
	"fmt"
	"testing"

	"gfx.cafe/util/go/gun"
)

var exampleConfigOne struct {
	Field1 string `yaml:"field_one"`
	Field2 int32
}

func TestConfigOne(t *testing.T) {

	gun.Load(&exampleConfigOne)
	fmt.Printf("t: %+v\n", exampleConfigOne)
}
