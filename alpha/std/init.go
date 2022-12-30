package std

func ReplaceIndex(str string, char rune, index int) string {
	bites := make([]rune, len(str))
	for i := 0; i < len(str); i++ {
		if i == index {
			bites[i] = char
		} else {
			bites[i] = rune(str[i])
		}
	}
	return string(bites)
}
