package validator

import "strconv"

func CheckLuhn(number string) bool {
	if len(number) == 0 {
		return false
	}

	var sum int

	parity := len(number) % 2
	for i := 0; i < len(number); i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
