// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package slices

// ContainsString determines whether a string is contained in a slice.
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
