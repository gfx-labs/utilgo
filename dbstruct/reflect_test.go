package dbstruct_test

import (
	"testing"

	"gfx.cafe/util/go/dbstruct"
	"github.com/stretchr/testify/require"
)

type User struct {
	Name     string `db:"name"`
	Status   string
	Password string `db:"-"`
	Id       int    `db:"id"`
}

func TestMakeQuery(t *testing.T) {
	query, err := dbstruct.MakeQuery("users", &User{})
	require.NoError(t, err)
	require.EqualValues(t, `insert into users ("name","status","id") values ($1,$2,$3) `, query)
}

func TestGetArgs(t *testing.T) {
	args, err := dbstruct.GetArgs(&User{
		"foo", "bar", "lorem", 5,
	})
	require.NoError(t, err)
	require.Len(t, args, 3)
	require.EqualValues(t, args[0], "foo")
	require.EqualValues(t, args[1], "bar")
	require.EqualValues(t, args[2], 5)
}

type IgnoreMe struct {
	IgnoreMe string
}
type DontIgnoreMe struct {
	DontIgnoreMe string
}
type SuperUser struct {
	Class    string
	User     User `db:",embedded"`
	IgnoreMe `db:"-"`
	DontIgnoreMe
}

func TestMakeQueryComplex(t *testing.T) {
	su := &SuperUser{
		Class:        "fuh",
		User:         User{"foo", "bar", "lorem", 5},
		IgnoreMe:     IgnoreMe{"ignore me"},
		DontIgnoreMe: DontIgnoreMe{"dont ignore me"},
	}
	query, err := dbstruct.MakeQuery("users", su)
	require.NoError(t, err)
	require.EqualValues(t, `insert into users ("class","name","status","id","dontignoreme") values ($1,$2,$3,$4,$5) `, query)
	args, err := dbstruct.GetArgs(su)
	require.NoError(t, err)
	require.Len(t, args, 5)
	require.EqualValues(t, args[0], "fuh")
	require.EqualValues(t, args[1], "foo")
	require.EqualValues(t, args[2], "bar")
	require.EqualValues(t, args[3], 5)
	require.EqualValues(t, args[4], "dont ignore me")
}
