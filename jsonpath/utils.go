package jsonpath

func UtilsJoin(delimiter string, warp string, stringValues []string) string {
	if len(stringValues) == 0 {
		return ""
	}
	r := ""
	for _, stringValue := range stringValues {
		r += delimiter + warp + stringValue
	}
	return r
}
