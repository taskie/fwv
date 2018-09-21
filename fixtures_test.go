package fwv

var records01 = [][]string{
	{"abc", "あいう", "αβγ", "abc", "あいう", "αβγ", "abc", "あいう", "αβγ"},
	{"abc", "あいう", "αβγ", "あいう", "αβγ", "abc", "αβγ", "abc", "あいう"},
	{"abc", "あいう", "αβγ", "abc", "あいう", "αβγ", "あいう", "αβγ", "abc"},
}

var records02 = [][]string{
	{"a", "b", ""},
	{"c", "", "d"},
	{"", "e", "f"},
	{"", "", ""},
}

var records03 = [][]string{
	{"あ", "いう", "えおか", "き"},
	{"アイ", "ウ", "エ", "オカキ"},
}

var fwv01 = `abc あいう αβγ abc あいう αβγ abc あいう αβγ
abc あいう αβγ あいう αβγ abc αβγ abc あいう
abc あいう αβγ abc あいう αβγ あいう αβγ abc
`

var fwvUseWidth01 = `abc あいう αβγ abc    あいう αβγ abc    あいう αβγ
abc あいう αβγ あいう αβγ abc    αβγ abc    あいう
abc あいう αβγ abc    あいう αβγ あいう αβγ abc
`

var fwvUseWidthEaaHalf01 = `abc あいう αβγ abc    あいう αβγ abc    あいう αβγ
abc あいう αβγ あいう αβγ    abc αβγ    abc    あいう
abc あいう αβγ abc    あいう αβγ あいう αβγ    abc
`

var fwv02 = `a b
c   d
  e f

`

var fwvUseWidth03 = `あ   いう えおか き
アイ ウ   エ     オカキ
`

var fwvUseWidthDelimited03 = `あ  |いう|えおか|き
アイ|ウ  |エ    |オカキ
`

var csv01 = `abc,あいう,αβγ,abc,あいう,αβγ,abc,あいう,αβγ
abc,あいう,αβγ,あいう,αβγ,abc,αβγ,abc,あいう
abc,あいう,αβγ,abc,あいう,αβγ,あいう,αβγ,abc
`
