package utils

func InitializeMap() map[string]interface{} {
	mapData := make(map[string]interface{})
	mapData["Links"] = *Links

	return mapData
}
