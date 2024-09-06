package envvariable

type EnvVariableService interface {
	Load(filenames ...string) error
	Get(key string) string
}
