package utils

// Reverse Reverses the array
func Reverse(elements []string) []string {
	for i := len(elements)/2 - 1; i >= 0; i-- {
		opp := len(elements) - 1 - i
		elements[i], elements[opp] = elements[opp], elements[i]
	}
	return elements
}

// RemoveDuplicates Remove duplicate items from an array.
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
