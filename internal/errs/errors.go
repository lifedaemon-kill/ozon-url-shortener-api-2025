package errs

import "errors"

var (
	ErrorYamlConfigFileNotFound = errors.New("yaml config not found")
	ErrorEnvConfigFileNotFound  = errors.New("env config not found")
)
