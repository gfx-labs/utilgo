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
	gun.LoadFile("./config.yml", &exampleConfigOne)
	gun.LoadEnvVars(&exampleConfigOne)
	fmt.Printf("t: %+v\n", exampleConfigOne)
}
