package goenv

import "strings"

func contains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

func containsValue(ss []string, key string) (string, bool) {
	for i := range ss {
		prefix := key + "="
		if strings.HasPrefix(ss[i], prefix) {
			return strings.TrimPrefix(ss[i], prefix), true
		}
	}

	return "", false
}
