package parse

import "io"

// Parser is the interface that parses input data to some other structure.
// The intention is to isolate all the parsing logic details behind this simple
// interface. Parser should be safe for concurrent use by multiple goroutines.
type Parser interface {
	Parse(input io.Reader) (interface{}, error)
}
