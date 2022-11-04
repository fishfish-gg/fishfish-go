package fishfish

// Valid domain categories for the FishFish API
var validCategories = []string{"safe", "malware", "phishing"}

func validCategory(category string) bool {
	for _, v := range validCategories {
		if v == category {
			return true
		}
	}

	return false
}
