package util


func StringInSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}


func CheckDuplicate(arr []string) bool {
	seen := make(map[string]bool)
	for _, num := range arr {
		if seen[num] {
			return true
		}
		seen[num] = true
	}
	return false
}
