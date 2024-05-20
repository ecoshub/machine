package machine

const (
	Zero           string = "0"
	Space          string = " "
	Comma          string = ","
	Point          string = "."
	Column         string = ":"
	AlphaUpper     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaLower     string = "abcdefghijklmnopqrstuvwxyz"
	Alpha          string = AlphaLower + AlphaUpper
	PositiveDigits string = "123456789"
	Digits         string = Zero + PositiveDigits
	HexLower       string = Zero + PositiveDigits + "abcdef"
	HexUpper       string = Zero + PositiveDigits + "ABCDEF"
	Hex            string = HexLower + HexUpper
	Underline      string = "_"
	At             string = "@"
	Dolar          string = "$"
	AlphaNumerical string = AlphaUpper + AlphaLower + Digits
)

var (
	PredefinedTransitions = map[string][]rune{
		"_zero_":            []rune(Zero),
		"_space_":           []rune(Space),
		"_comma_":           []rune(Comma),
		"_point_":           []rune(Point),
		"_column_":          []rune(Column),
		"_alpha_upper_":     []rune(AlphaUpper),
		"_alpha_lower_":     []rune(AlphaLower),
		"_alpha_":           []rune(Alpha),
		"_positive_digits_": []rune(PositiveDigits),
		"_digits_":          []rune(Digits),
		"_hex_lower_":       []rune(HexLower),
		"_hex_upper_":       []rune(HexUpper),
		"_hex_":             []rune(Hex),
		"_underline_":       []rune(Underline),
		"_at_":              []rune(At),
		"_dolar_":           []rune(Dolar),
		"_alpha_numerical_": []rune(AlphaNumerical),
	}
)
