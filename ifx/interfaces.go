package ifx

import "io"

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

type EncoderFunc func(io.Writer) Encoder
type DecoderFunc func(io.Reader) Decoder

type Compressor func(w io.Writer) (io.WriteCloser, error)

type Decompressor func(r io.Reader) io.ReadCloser
