package utiles

// StrListContain check if an string list contain an item
func StrListContain(list []string, term string) bool {
	for _, i := range list {
		if i == term {
			return true
		}
	}
	return false
}
