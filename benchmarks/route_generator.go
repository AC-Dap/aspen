package benchmarks

import (
	"math/rand"
	"strings"
)

const SEGMENT_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// GenerateRandomPath generates a random path for testing purposes.
// It varies the number of segments in the path and the length of each segment.
func GenerateRandomPath(rng *rand.Rand) string {
	// Can have 1 - 5 num_segments, weighted around 3 num_segments.
	// 	1: 5%, 2: 25%, 3: 40%, 4: 25%, 5: 5%
	var num_segments int
	segments_rng := rng.Float64()
	if segments_rng < 0.05 {
		num_segments = 1
	} else if segments_rng < 0.3 {
		num_segments = 2
	} else if segments_rng < 0.7 {
		num_segments = 3
	} else if segments_rng < 0.95 {
		num_segments = 4
	} else {
		num_segments = 5
	}

	// Generate each segment with a length of 1 - 10 characters, weighted around 5 characters.
	var sb strings.Builder
	for i := 0; i < num_segments; i++ {
		// Every segment is prepended with a slash, including the first segment.
		sb.WriteByte('/')

		segment_length := int(rng.NormFloat64()*1.5 + 5.5)
		if segment_length < 1 {
			segment_length = 1
		} else if segment_length > 10 {
			segment_length = 10
		}

		for j := 0; j < segment_length; j++ {
			// Randomly choose a character from SEGMENT_CHARS
			char := SEGMENT_CHARS[rng.Intn(len(SEGMENT_CHARS))]
			sb.WriteByte(char)
		}
	}

	return sb.String()
}

// GenerateRandomPaths generates a slice of random paths.
func GenerateRandomPaths(rng *rand.Rand, num_paths int) []string {
	paths := make([]string, num_paths)
	for i := range num_paths {
		paths[i] = GenerateRandomPath(rng)
	}
	return paths
}
