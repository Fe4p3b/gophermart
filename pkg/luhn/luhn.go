package luhn

func OnlyDigits(b []byte) bool {
	for _, n := range b {
		if n >= 48 && n <= 57 {
			continue
		}
		return false
	}
	return true
}

func Luhn(s []byte) bool {
	sum := 0
	for i := len(s) - 1; i >= 0; i-- {
		d := int(s[i] - '0')

		if i%2 == 1 {
			d = d * 2
		}
		sum += d / 10
		sum += d % 10

	}
	return sum%10 == 0
}
