# gun

gun is a very opinionated and simple config library


## how to add


```
go get gfx.cafe/open/gun
```

## how to use
simply create a struct that is json encodable into what you need

struct keys must be alphanumeric + _

json is a subset of yaml, so your config file can be in JSON or YAML, however use yaml as your tags

config files key by default is all lowercase of the field (this is different from JSON!)

env vars WILL ALWAYS override, and ALWAYS ALL UPPERCASE

does NOT work with nested variables.

env vars ONLY works with integer types, float types, and string. all other types will PANIC!!!


```
package main


var Config struct {
  Field string `yaml:"some_field"`
  Crack int
  Arango_USER dog
}

// the env vars for this struct are
// SOME_FIELD, CRACK, ARANGO_USER
// the yaml/json for this struct are
// field, crack, arango_user


func init() {
  gun.Load(&Config)
}

func main() {
  _ = Config.Field
}
```

It will look for the config file in the following places.

1. env var $GUN_CONFIG_FILE
2. /config/config.yml
3. /config.yml
5. /config/config.yaml
6. /config.yaml
5. /config/config.json
6. /config.json


## important

gun will read config file only once, when you ask it to read.

this means that if env variables/configs change, gun will not pick it up!
