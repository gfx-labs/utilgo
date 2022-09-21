package gun

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(i any) {
	if errs := mustLoad(i); len(errs) > 0 {
		for idx, v := range errs {
			fmt.Printf("round %d: %s\n", idx, v)
		}
		panic("no suitable configs!")
	}
	if err := LoadEnvVars(i); err != nil {
		panic(err)
	}
}

func mustLoad(i any) []error {
	dir, _ := os.UserHomeDir()
	errs := []error{}
	files := make([]string, 0, 4)
	if gcf := os.Getenv("GUN_CONFIG_FILE"); gcf != "" {
		files = append(files, gcf)
	}
	files = append(files, []string{
		"/config/config.yml", "/config.yml", path.Join(dir, ".gfx/config.yml"),
		"/config/config.yaml", "/config.yaml", path.Join(dir, ".gfx/config.yaml"),
		"/config/config.json", "/config.json", path.Join(dir, ".gfx/config.json"),
	}...)
	for _, v := range files {
		err := LoadFile(v, i)
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}
	return nil
}

func LoadFile(file string, i any) error {
	bts, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bts, i)
	if err != nil {
		return err
	}
	return nil
}

func LoadEnvVars(i any) error {
	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("%v is not a pointer type", i)
	}
	rv = rv.Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {

		envVar, ok := rt.Field(i).Tag.Lookup("yaml")
		if ok {
			envVar = strings.Split(envVar, ",")[0]
		} else {
			envVar, ok = rt.Field(i).Tag.Lookup("gun")
			if ok {
				envVar = strings.Split(envVar, ",")[0]
			} else {
				envVar = strings.ToUpper(rt.Field(i).Name)
				log.Println(rt.Field(i).Name, envVar)
			}
		}
		evs := os.Getenv(strings.ToUpper(envVar))
		if evs != "" {
			switch t := rv.Field(i).Interface().(type) {
			case int, int8, int16, int32, int64:
				rv.Field(i).SetInt(int64(mustEnvInt(envVar)))
			case uint, uint8, uint16, uint32, uint64:
				rv.Field(i).SetUint(uint64(mustEnvInt(envVar)))
			case string:
				rv.Field(i).SetString(evs)
			case []string:
				rv.Field(i).Set(reflect.ValueOf(strings.Split(evs, ",")))
			case []int:
				a := make([]int, 0)
				for _, v := range strings.Split(evs, ",") {
					a = append(a, mustInt(v))
				}
				rv.Field(i).Set(reflect.ValueOf(a))
			case []byte:
				bts, err := hex.DecodeString(evs)
				if err == nil {
					rv.Field(i).SetBytes(bts)
				}
			default:
				return fmt.Errorf("unsupported type %v", t)
			}
		}
	}

	return nil
}
func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
func mustEnvInt(s string) int {
	i, err := strconv.Atoi(os.Getenv(s))
	if err != nil {
		panic(err)
	}
	return i
}
