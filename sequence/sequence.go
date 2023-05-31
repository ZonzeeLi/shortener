package sequence

// Sequence is an interface for generating sequence numbers.
type Sequence interface {
	Next() (uint64, error)
}
