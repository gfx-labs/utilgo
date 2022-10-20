package rangerunner

import (
	"database/sql/driver"

	"github.com/RoaringBitmap/roaring/roaring64"
)

type Seq struct {
	max int
	m   *roaring64.Bitmap
}

func NewSeq(n int) *Seq {
	return &Seq{
		max: n,
		m:   roaring64.New(),
	}
}

func (s *Seq) Load(r *roaring64.Bitmap) {
	s.m.Or(r)
}

// sets max of the requested sequence
func (s *Seq) SetMax(n int) {
	s.max = n
}

// adds range from start to stop inclusive
func (s *Seq) Add(start, stop int) {
	s.m.AddRange(uint64(start), uint64(stop+1))
}

func (s *Seq) Remove(start, stop int) {
	s.m.RemoveRange(uint64(start), uint64(stop+1))
}

func (s *Seq) U() *roaring64.Bitmap {
	return s.m
}

func (s *Seq) Next(size int) [2]int {
	iter := s.m.Iterator()
	var val uint64
	for iter.HasNext() {
		val = iter.Next()
	}
	return [2]int{int(val), int(val) + size - 1}
}

func (s *Seq) Str() string {
	ss, _ := s.m.ToBase64()
	return ss
}

func (s *Seq) Value() (driver.Value, error) {
	return s.Str(), nil
}

func (s *Seq) Scan(src any) error {
	switch v := src.(type) {
	case string:
		s.m.FromBase64(v)
	}
	return nil
}

func (s *Seq) UnmarshalJSON(xs []byte) error {
	_, err := s.m.FromBase64(string(xs))
	if err != nil {
		return err
	}
	return nil
}

func (s *Seq) MarshalJSON() ([]byte, error) {
	return []byte(s.Str()), nil
}
