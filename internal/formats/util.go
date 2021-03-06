package formats

import "github.com/arykalin/format-bot/internal/formats/data_getter"

func contains(s []data_getter.Tag, str data_getter.Tag) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// subset returns true if the first array is completely
// contained in the second array. There must be at least
// the same number of duplicate values in second as there
// are in first.
func subset(first, second []data_getter.Tag) bool {
	for _, s := range first {
		if !contains(second, s) {
			return false
		}
	}
	return true
}
