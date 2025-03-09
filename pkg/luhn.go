package pkg

// ValidateLuhn проверяет, соответствует ли строка чисел алгоритму Луна (Luhn Algorithm).
// Алгоритм Луна используется для проверки корректности чисел, таких как номера кредитных карт.
//
// Аргументы:
//   - number: строка, представляющая число для проверки, например, номер кредитной карты.
//
// Возвращает:
//   - true, если число соответствует алгоритму Луна.
//   - false, если число не соответствует алгоритму Луна.
func ValidateLuhn(number string) bool {
	if len(number) == 0 {
		return false
	}

	sum := 0
	isEven := false

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if isEven {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}
