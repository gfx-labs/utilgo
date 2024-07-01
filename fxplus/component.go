package componentname

type ComponentName string

func WithName(x string) func() ComponentName {
	return func() ComponentName {
		return ComponentName(x)
	}
}
