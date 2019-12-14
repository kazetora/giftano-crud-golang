package helpers

func Difference(a, b []int) (diff []int) {
	ma := make(map[int]bool)
	mb := make(map[int]bool)

	for _, item := range a {
		ma[item] = false
	}

	for _, item := range b {
			mb[item] = true
	}

	for _, item := range a {
			if _, ok := mb[item]; !ok {
				if !ma[item] {
					diff = append(diff, item)
					ma[item] = true
				}
			}
	}
	return diff
}