package util

type Encoder interface {
	Encode(message string) string
}

type Decoder interface {
	Decode(encodedMessage string) string
}
