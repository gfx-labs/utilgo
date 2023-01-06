# gun

gun is a very opinionated and simple config library


## how to add


```
go get gfx.cafe/open/gun
```

## how to use
simply create a struct that is json encodable into what you need

struct keys must be alphanumeric + _

your config file can be .yml .yaml .json or .env (.env is only read if it is in the current directory)

env vars WILL ALWAYS override

does NOT work with nested variables.

```
package main

var exampleConfigOne struct {
	Field1     string `yaml:"field_one" env:"FIELD_ONE" json:"f_1" default:"IM DEFAULT HI"`
	Field2     int32
	MANY_FIELD []string
}

// the env vars for this struct are
// FIELD_ONE, Field2, MANY_FIELD
// the yaml for this struct are
// field_one, Field2, MANY_FIELD
// the json for this struct are
// f_1, Field2, MANY_FIELD

func init() {
  gun.Load(&Config)
}

func main() {
  _ = Config.Field1
}
```

You can also call

```gun.LoadPrefix(&Config, "prefix")```

Unless a field is tagged otherwise gun will look for exactly the name of that field
to look for a different name use yaml:"" env:"" or json:"" tags

to set a default value for that config item use the default:"" tag,
this is overridden by any file or ENV present

to load ENV variables prefixed with "prefix"
and all prefix.yml prefix.yaml prefix.json files in the same locations listed below
this will not apply to values inside of .env files

It will look for the config file in the following places.
Gun loads the config in the order listed, the lower items will override the higher.
For example a value set in ~/.gfx/config.yml will override the same value set in /config.yml

1. /config.yml
2. /config.yaml
3. /config.json
4. /config/config.yml
5. /config/config.yaml
6. /config/config.json
7. ~/.gfx/config.yml
8. ~/.gfx/config.yaml
9. ~/.gfx/config.json
10. .env

## important

gun will read config file only once, when you ask it to read.

this means that if env variables/configs change, gun will not pick it up!
