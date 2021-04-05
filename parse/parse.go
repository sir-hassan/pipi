package parse

import "io"

type Parser interface {
	Parse(input io.Reader) (interface{}, error)
}
