package randutils

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) { pseudo.Shuffle(n, swap) }
