package config

type EnvConfig struct {
	COMPUTING_POWER         int
	TIME_ADDITION_MS        int
	TIME_SUBSTRACTION_MS    int
	TIME_DIVISIONS_MS       int
	TIME_MULTIPLICATIONS_MS int
}

func DefaultEnvConfig() *EnvConfig {

	return &EnvConfig{
		COMPUTING_POWER:         6,
		TIME_ADDITION_MS:        15,
		TIME_SUBSTRACTION_MS:    15,
		TIME_DIVISIONS_MS:       20,
		TIME_MULTIPLICATIONS_MS: 20,
	}
}

type ErrConfig struct {
	ELEMENTNOTFOUND int
	EMPTYSLICE      int
}

func DefaultErrConfig() *ErrConfig {
	return &ErrConfig{
		ELEMENTNOTFOUND: -1,
		EMPTYSLICE:      -2,
	}
}

type HttpConfig struct {
	Port string
}

func DeafultHttpConfig() *HttpConfig {
	return &HttpConfig{
		Port: ":8080",
	}
}

type GrpcConfig struct {
	URL              string
	InternalServ     string
	TaskServ         string
	ProcessedExpServ string
}

func DefaultGrpcConfig() *GrpcConfig {
	return &GrpcConfig{
		URL:              "http://localhost:8080/internal/task",
		InternalServ:     "localhost:5000",
		TaskServ:         "localhost:5001",
		ProcessedExpServ: "localhost:5002",
	}
}

type DBConfig struct {
	CnnectionRefused int
	PrepareFailed    int
	ScanFailed       int
	NoRows           int
}

func DefaultDBConfig() *DBConfig {
	return &DBConfig{
		CnnectionRefused: -1,
		PrepareFailed:    -2,
		ScanFailed:       -3,
		NoRows:           -4,
	}
}

type JWTConfig struct {
	PasswordSigningJwt string
}

func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		PasswordSigningJwt: "2zTvQ4h7XQy4kL0w9J9n8K7m6p5r3s2t1u0vWxYzAeBc",
	}
}

type Config struct {
	*EnvConfig
	*ErrConfig
	*HttpConfig
	*GrpcConfig
	*DBConfig
	*JWTConfig
}

func DefaultConfig() *Config {
	return &Config{
		DefaultEnvConfig(),
		DefaultErrConfig(),
		DeafultHttpConfig(),
		DefaultGrpcConfig(),
		DefaultDBConfig(),
		DefaultJWTConfig(),
	}
}
