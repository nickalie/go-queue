package queue

func increaseString(value string) string {
	if value == "" {
		value = "a"
	} else {
		l := len(value)
		r := value[l-1]

		if r < 122 {
			r++
			value = string(append([]byte(value)[:l-1], r))
		} else {
			value = string(append([]byte(value), 97))
		}
	}

	return value
}
