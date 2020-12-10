package ot

type TextEncodingType int

const (
	TextEncodingTypeUTF8 = iota
	TextEncodingTypeUTF16
)

// Use utf-8 by default
var TextEncoding = TextEncodingTypeUTF8
