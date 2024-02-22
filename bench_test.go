package machine_test

import (
	"machine"
	"regexp"
	"testing"
)

var (
	validValue string = "-12311098123873"
	reg        *regexp.Regexp
	sm         *machine.StateMachine
)

func init() {
	var err error
	reg, err = regexp.Compile("^(0|-*[1-9]+[0-9]*)$")
	if err != nil {
		panic(err)
	}

	sm, err = machine.Build([]byte(`

$zero: '0'
$minus: '-'
$natural_number: "123456789"
$digits: {$zero} + {$natural_number}

@positive_integer: [
	_start_ > {$natural_number} > number > {$digits} > number > * > _end_
]

@negative_integer: [
	_start_ > {$minus} > {s1: @positive_integer} > * > _end_
]

@integer: [
	_start_ > * > {s2: @positive_integer} > * > _end_,
	_start_ > * > {s3: @negative_integer} > * > _end_
]

@_main_: [
	_start_ > * > {s4: @integer} > * > _end_
]

`))
	if err != nil {
		panic(err)
	}
}

func BenchmarkValidRegex(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reg.Match([]byte(validValue))
	}
}

func BenchmarkValidSM(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sm.Match(validValue)
	}
}
