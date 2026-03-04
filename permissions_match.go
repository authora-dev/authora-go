package authora

import "strings"

func matchSegment(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, pattern[:len(pattern)-1])
	}
	return pattern == value
}

func MatchPermission(pattern, resource string) bool {
	p := strings.Split(pattern, ":")
	r := strings.Split(resource, ":")
	if len(p) != len(r) {
		return false
	}
	for i := range p {
		if !matchSegment(p[i], r[i]) {
			return false
		}
	}
	return true
}

func MatchAnyPermission(patterns []string, resource string) bool {
	for _, p := range patterns {
		if MatchPermission(p, resource) {
			return true
		}
	}
	return false
}
