package enc

import "io"

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

type EncoderFunc func(io.Writer) Encoder
type DecoderFunc func(io.Reader) Decoder
