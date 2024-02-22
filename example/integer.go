package main

import (
	"fmt"
	"machine"
)

func main() {
	sm, err := machine.Build([]byte(`

$zero: '0'
$minus: '-'
$natural_number: "123456789"
$digits: {$zero} + {$natural_number}

@positive_integer: [
	_start_ > {$natural_number} > number > {$digits} > number > * > _end_
]

@negative_integer: [
	_start_ > {$minus} > {@positive_integer} > * > _end_
]

@integer: [
	_start_ > * > {@positive_integer} > * > _end_,
	_start_ > * > {@negative_integer} > * > _end_
]

@_main_: [
	_start_ > * > {@integer} > * > _end_
]

`))
	if err != nil {
		fmt.Println(err)
		return
	}

	text := "-31123981402381"
	err = sm.Match(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("OK")
}
