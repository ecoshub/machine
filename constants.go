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
	AlphaNumeric   string = AlphaUpper + AlphaLower + Digits
)
