package repository

type EnvVar interface {
	bool | int | int32 | int64 | float32 | float64
}

// Repository is a interface tha has func to get the value of a config
type Repository interface {
	Get(string) string
}
