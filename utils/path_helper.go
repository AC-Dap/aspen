package utils

import (
	"strings"

	"github.com/julienschmidt/httprouter"
)

const (
	FIXED_SEGMENT PathSegmentType = iota
	NAMED_SEGMENT
	CATCH_ALL_SEGMENT
)

type PathSegmentType int

type PathSegment struct {
	Type  PathSegmentType
	Value string
}

type Path []PathSegment

// ParsePath takes a path string and returns a Path, which is a slice of PathSegment's.
func ParsePath(path string) Path {
	// Split the path by '/' and parse each segment
	segments := strings.Split(path, "/")
	parsedSegments := make(Path, len(segments))

	for i, segment := range segments {
		var parsedSegment PathSegment

		if len(segment) == 0 {
			parsedSegment = PathSegment{
				Type:  FIXED_SEGMENT,
				Value: "",
			}
		} else {
			// Determine the type of segment based on the first character
			switch segment[0] {
			case ':': // Named segment
				parsedSegment = PathSegment{
					Type:  NAMED_SEGMENT,
					Value: segment[1:], // Exclude the leading ':'
				}
			case '*': // Catch-all segment
				parsedSegment = PathSegment{
					Type:  CATCH_ALL_SEGMENT,
					Value: segment[1:], // Exclude the leading '*'
				}
			default: // Fixed segment
				parsedSegment = PathSegment{
					Type:  FIXED_SEGMENT,
					Value: segment,
				}
			}
		}

		parsedSegments[i] = parsedSegment
	}

	return parsedSegments
}

// IsProxyCompatible checks if the given proxy path has the same variables as this path.
func (p Path) IsProxyCompatible(proxy Path) bool {
	// Check that each variable in the proxy has a corresponding variable in the given path, and vice versa.
	// This means that we have a 1-to-1 mapping of variables.

	for _, pSegment := range proxy {
		if pSegment.Type == FIXED_SEGMENT {
			continue // Fixed segments can match anything
		}

		found := false
		for _, rSegment := range p {
			if pSegment.Type == rSegment.Type && pSegment.Value == rSegment.Value {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	for _, rSegment := range p {
		if rSegment.Type == FIXED_SEGMENT {
			continue // Fixed segments can match anything
		}

		found := false
		for _, pSegment := range proxy {
			if rSegment.Type == pSegment.Type && rSegment.Value == pSegment.Value {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// ConstructPath takes a set of variables and constructs the path string, replacing named segments with their values.
func (p Path) ConstructPath(variables httprouter.Params) string {
	var sb strings.Builder
	for i, segment := range p {
		// Separate segments with a '/'
		if i > 0 {
			sb.WriteString("/")
		}

		switch segment.Type {
		case FIXED_SEGMENT:
			// Fixed segments are written as-is
			sb.WriteString(segment.Value)

		case NAMED_SEGMENT:
			// Look up variable value and substitute that in
			sb.WriteString(variables.ByName(segment.Value))

		case CATCH_ALL_SEGMENT:
			// Catch-all segments capture the leading '/', so we need to ensure we don't add an extra one
			sb.WriteString(variables.ByName(segment.Value)[1:])
		}
	}

	return sb.String()
}

func (p Path) String() string {
	// Join the segments into a string, using '/' as the separator
	var sb strings.Builder
	for i, segment := range p {
		if i > 0 {
			sb.WriteString("/")
		}
		switch segment.Type {
		case FIXED_SEGMENT:
			sb.WriteString(segment.Value)
		case NAMED_SEGMENT:
			sb.WriteString(":" + segment.Value)
		case CATCH_ALL_SEGMENT:
			sb.WriteString("*" + segment.Value)
		}
	}
	return sb.String()
}
