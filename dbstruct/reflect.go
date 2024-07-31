package dbstruct

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/structtag"
)

var driverType = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

type bindingCache struct {
	m  map[reflect.Type]*structTableBinding
	mu sync.Mutex
}

var globalBindingCache = &bindingCache{m: make(map[reflect.Type]*structTableBinding)}

type structTableBinding struct {
	name   string
	t      reflect.Type
	fields []structField
}

type structField struct {
	// type of field
	t reflect.Type
	// index of field
	f []int
	// column name of fields
	c string
}

func GetArgs(item any) ([]any, error) {
	value, err := getStructValue(item)
	if err != nil {
		return nil, err
	}

	// grab type binding
	binding, err := calculateTypeBinding(item)
	if err != nil {
		return nil, err
	}

	values := make([]any, 0, len(binding.fields))
	for _, v := range binding.fields {
		f := value.FieldByIndex(v.f)
		values = append(values, f.Interface())
	}
	return values, nil
}

func MakeQuery(table string, item any) (string, error) {
	// grab type binding
	binding, err := calculateTypeBinding(item)
	if err != nil {
		return "", err
	}
	if len(binding.fields) == 0 {
		return "", fmt.Errorf("struct must have at least one field to insert")
	}

	bb := new(strings.Builder)
	pb := new(strings.Builder)
	for idx, v := range binding.fields {
		bb.WriteString(`"` + v.c + `"`)
		pb.WriteString("$")
		pb.WriteString(strconv.Itoa(idx + 1))
		if idx+1 != len(binding.fields) {
			bb.WriteString(",")
			pb.WriteString(",")
		}
	}

	// form the query not
	query := `insert into ` + table +
		" (" + bb.String() + `) ` +
		`values (` + pb.String() + `) `
	return query, nil
}

func calculateTypeBinding(item any) (*structTableBinding, error) {
	// evaluate the struct
	value, err := getStructValue(item)
	if err != nil {
		return nil, err
	}
	// key by type name
	typ := value.Type()

	globalBindingCache.mu.Lock()
	defer globalBindingCache.mu.Unlock()
	typeBinding, ok := globalBindingCache.m[typ]
	if ok {
		return typeBinding, nil
	}
	fields, err := fieldToInfo(typ, nil)
	if err != nil {
		return nil, err
	}
	typeBinding = &structTableBinding{
		name:   typ.Name(),
		t:      typ,
		fields: fields,
	}
	globalBindingCache.m[typ] = typeBinding
	return typeBinding, nil
}

func fieldToInfo(typ reflect.Type, index []int) ([]structField, error) {
	typ, err := getStructType(typ)
	if err != nil {
		return nil, err
	}
	output := make([]structField, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// skip private fields
		if !field.IsExported() {
			continue
		}
		idx := append([]int{}, append(index, i)...)

		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, err
		}
		columnName := strings.ToLower(field.Name)
		dbTag, err := tags.Get("db")
		if err == nil {
			columnName = dbTag.Name
		}
		if columnName == "-" {
			continue
		}
		var isValuer bool
		if field.Type.Implements(driverType) || field.Type.Implements(scannerType) {
			isValuer = true
		}
		// if anonymous field or embedded field, we need to iterate over it
		if !isValuer && (field.Anonymous || (err == nil && dbTag != nil && dbTag.HasOption("embedded"))) {
			fields, err := fieldToInfo(field.Type, idx)
			if err != nil {
				return nil, err
			}
			output = append(output, fields...)
			continue
		}
		// otherwise, assume its a field.
		output = append(output, structField{
			t: field.Type,
			f: idx,
			c: columnName,
		})
	}
	return output, nil
}

func getTypeBinding(item any) (*structTableBinding, error) {
	structValue, err := getStructValue(item)
	if err != nil {
		return nil, err
	}
	binding, err := calculateTypeBinding(structValue)
	if err != nil {
		return nil, err
	}
	return binding, nil
}

func getStructValue(item any) (reflect.Value, error) {
	value := reflect.ValueOf(item)
	switch value.Kind() {
	case reflect.Struct:
		return value, nil
	case reflect.Pointer:
		value = value.Elem()
		if value.Kind() == reflect.Struct {
			return value, nil
		}
	}
	return value, fmt.Errorf("can only work on struct or struct pointer. got %s", value.Kind())
}
func getStructType(typ reflect.Type) (reflect.Type, error) {
	switch typ.Kind() {
	case reflect.Struct:
		return typ, nil
	case reflect.Pointer:
		typ = typ.Elem()
		if typ.Kind() == reflect.Struct {
			return typ, nil
		}
	}
	return typ, fmt.Errorf("can only work on struct or struct pointer. got %s", typ.Kind())
}
