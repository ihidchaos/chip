package compare

// Eqer can be used to determine if this value is equal to the other.
// The semantics of equals is that the two value are interchangeable
// in the Hugo templates.
type Eqer interface {
	// Eq returns whether this value is equal to the other.
	// This is for internal use.
	Eq(other any) bool
}

// ProbablyEqer is an equal check that may return false positives, but never
// a false negative.
type ProbablyEqer interface {
	// For internal use.
	ProbablyEq(other any) bool
}

// Comparer can be used to compare two values.
// This will be used when using the le, ge etc. operators in the templates.
// Compare returns -1 if the given version is less than, 0 if equal and 1 if greater than
// the running version.
type Comparer interface {
	Compare(other any) int
}
