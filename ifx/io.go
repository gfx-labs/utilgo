package ifx

import "io"

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

type StreamWriter func(io.Writer) Encoder
type StreamReader func(io.Reader) Decoder

type Compressor func(w io.Writer) (io.WriteCloser, error)
type Decompressor func(r io.Reader) io.ReadCloser
