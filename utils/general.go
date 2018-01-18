package utils

// Reverse Reverses the array
func Reverse(elements []string) []string {
	for i := len(elements)/2 - 1; i >= 0; i-- {
		opp := len(elements) - 1 - i
		elements[i], elements[opp] = elements[opp], elements[i]
	}
	return elements
}
