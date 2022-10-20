package ifx

import "io"

// implements fn Flush() error
type Flusher interface {
	Flush() error
}

// implements fn Encode(any) error
type Encoder interface {
	Encode(any) error
}

// an encoder and io closer
type EncodeCloser interface {
	Encoder
	io.Closer
}

// encoder and flusher
type EncodeFlusher struct {
	Encoder
	Flusher
}

// implements fn Decode(any) error
type Decoder interface {
	Decode(any) error
}

// decoder and io closer
type DecodeCloser interface {
	Decoder
	io.Closer
}

// decoder and flusher
type DecodeFlusher interface {
	Decoder
	Flusher
}

type StreamEncoder func(io.Writer) Encoder
type StreamDecoder func(io.Reader) Decoder

type StreamEncodeFlusher func(io.Writer) EncodeFlusher
type StreamDecodeFlusher func(io.Reader) DecodeFlusher

type StreamEncodeCloser func(io.Writer) EncodeCloser
type StreamDecodeCloser func(io.Reader) DecodeCloser
