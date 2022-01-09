package arrutils

func StringSliceHas(slice []string, searchedString string) bool {
	for _, str := range slice {
		if str == searchedString {
			return true
		}
	}
	return false
}
