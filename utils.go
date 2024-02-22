package machine

import "math/rand"

const (
	hex = "0123456789abcdef"
)

const (
	randomIDLength int = 8
)

func randomID() string {
	arr := make([]byte, randomIDLength)
	for i := range arr {
		index := rand.Intn(len(hex))
		arr[i] = hex[index]
	}
	return string(arr)
}

func resolveOptionalString(optionalString ...string) string {
	val := ""
	if len(optionalString) > 0 {
		val = optionalString[0]
	} else {
		val = randomID()
	}
	return val
}

func printCharString(values []rune) string {
	s := ""
	for _, v := range values {
		sc, _ := printChar(v)
		s += sc
	}
	return s
}

func printChar(char rune) (string, int) {
	stringChar := ""
	intChar := 0
	if char == '\n' || char == '\t' {
		if char == '\n' {
			stringChar = `\n`
			intChar = int('\n')
		}
		if char == '\t' {
			stringChar = `\t`
			intChar = int('\t')
		}
	} else {
		stringChar = string(char)
		intChar = int(char)
	}
	return stringChar, intChar
}

func rawStringToRunes(values string) []rune {
	arr := make([]rune, 0, len(values))
	escape := false
	for _, v := range values {
		if v == '\\' {
			if !escape {
				escape = true
				continue
			}
		}
		if escape {
			switch v {
			case 'n':
				arr = append(arr, rune('\n'))
			case 't':
				arr = append(arr, rune('\t'))
			case 'r':
				arr = append(arr, rune('\r'))
			case 'b':
				arr = append(arr, rune('\b'))
			default:
				arr = append(arr, v)
			}
			escape = false
			continue
		}
		arr = append(arr, v)
	}
	return arr
}
