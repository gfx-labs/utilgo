package ifx

import "io"

type Encoder interface {
	Encode(any) error
}

type Decoder interface {
	Decode(any) error
}

type StreamEncoder func(io.Writer) Encoder
type StreamDecoder func(io.Reader) Decoder

type Compressor func(w io.Writer) (io.WriteCloser, error)
type MustCompressor func(w io.Writer) io.WriteCloser

type Decompressor func(r io.Reader) (io.ReadCloser, error)
type MustDecompressor func(r io.Reader) io.ReadCloser
