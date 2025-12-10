package handlers

import (
	_ "embed"
)

//go:embed openapi.yml
var OpenAPISpec []byte
