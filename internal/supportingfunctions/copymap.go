package supportingfunctions

func CopyMap[K comparable, T any](oldMap map[K]T) map[K]T {
	newMap := make(map[K]T, len(oldMap))
	for key, value := range oldMap {
		newMap[key] = value
	}

	return newMap
}
