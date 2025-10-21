package weatherutil

var (
	cityIdMap = map[string]string{
		"1": "札幌", "13": "東京", "23": "名古屋", "27": "大阪", "40": "博多",
	}
	weatherIdMap = map[string]string{
		"12": "雨", "4": "くもり", "2": "晴れ",
	}
)

func IdToCityName(id string) string {
	return cityIdMap[id]
}

func WeatherIdToName(id string) string {
	return weatherIdMap[id]
}
