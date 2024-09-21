package matcher

// This is the map function, it takes a Match and the data segment. This function will
// be invoken when writing down to the file.
type MatchMapper func(md Match, segment []byte) ([]byte, bool)

// Performs a modification action.
type Modifier interface {
	// Replace the matched section with the given map function
	Replace(MatchMapper) (bool, error)

	// Insert the matched section before/after with the map function.
	InsertBefore(MatchMapper) (bool, error)
	InsertAfter(MatchMapper) (bool, error)

	// Remove that chuck of the mathing
	Remove() (bool, error)
}
