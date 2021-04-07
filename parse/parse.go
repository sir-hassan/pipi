package parse

import "io"

// Parser is the interface that parses input data to something other structure.
// The intention is to isolate all the parsing logic details behind this simple interface.
type Parser interface {
	Parse(input io.Reader) (interface{}, error)
}
