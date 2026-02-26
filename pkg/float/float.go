package float

import "strings"

func ParseFloatBase(value string, base int) float64 {
	value = strings.TrimSpace(value)
	negative := false
	switch value[0] {
	case '-':
		negative = true
		value = value[1:]
	case '+':
		value = value[1:]
	}

	parts := strings.SplitN(value, ".", 2)
	intPart := parts[0]
	var fracPart string
	if len(parts) == 2 {
		fracPart = parts[1]
	}

	var intValue float64
	for _, r := range intPart {
		d := digitValue(r)
		intValue = intValue*float64(base) + float64(d)
	}

	var fracValue float64
	if fracPart != "" {
		divisor := float64(base)
		for _, r := range fracPart {
			d := digitValue(r)

			fracValue += float64(d) / divisor
			divisor *= float64(base)
		}
	}

	result := intValue + fracValue
	if negative {
		result = -result
	}
	return result
}

func FormatFloatBase(x float64, base int) string {
	negative := x < 0
	if negative {
		x = -x
	}

	intPart := int64(x)
	fracPart := x - float64(intPart)

	digits := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	var intStr string
	if intPart == 0 {
		intStr = "0"
	} else {
		var intBuilder []byte
		for intPart > 0 {
			intBuilder = append([]byte{digits[intPart%int64(base)]}, intBuilder...)
			intPart /= int64(base)
		}
		intStr = string(intBuilder)
	}

	var fracStr string
	if fracPart > 0 {
		var fracBuilder []byte
		maxDigits := 17
		for i := 0; fracPart > 0 && i < maxDigits; i++ {
			fracPart *= float64(base)
			d := int(fracPart)
			fracBuilder = append(fracBuilder, digits[d])
			fracPart -= float64(d)
		}
		fracStr = "." + string(fracBuilder)
	}

	result := intStr + fracStr
	if negative {
		result = "-" + result
	}
	return result
}

func digitValue(r rune) int {
	switch {
	case '0' <= r && r <= '9':
		return int(r - '0')
	case 'a' <= r && r <= 'z':
		return int(r-'a') + 10
	case 'A' <= r && r <= 'Z':
		return int(r-'A') + 10
	default:
		return 0
	}
}
