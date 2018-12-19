package utils

//ContainsInt return true if int i exists in s
func ContainsInt(s []int, i int) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}
