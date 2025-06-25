package benchmarks

import "math/rand"

const RNG_SEED int64 = 0x1234567890abcdef

// GetRNG initializes and returns a random number generator with a fixed seed for reproducibility.
func GetRNG() *rand.Rand {
	return rand.New(rand.NewSource(RNG_SEED))
}
