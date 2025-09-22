package errs

import "errors"

var (
	ErrorConfigFileNotFound = errors.New("config file not found")

	ErrorRepositoryUrlEmpty  = errors.New("url is empty")
	ErrorRepositoryDuplicate = errors.New("url already exists")

	ErrorUrlServiceLinkNotFound = errors.New("link not found")
	ErrorUrlServiceInternal     = errors.New("internal error")

	ErrorAlreadyExist = errors.New("url already exist")
)
