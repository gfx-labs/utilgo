package fxplus

type ComponentName string

func Component(x string) func() ComponentName {
	return func() ComponentName {
		return ComponentName(x)
	}
}
