package huffman

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCalculateFrequencies(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected [256]uint64
	}{
		{
			name:  "normal",
			input: []byte("banana"),
			expected: [256]uint64{
				'a': 3,
				'b': 1,
				'n': 2,
			},
		},
		{
			name: "empty",
		},
		{
			name:  "repeated",
			input: []byte("aaaaaa"),
			expected: [256]uint64{
				'a': 6,
			},
		},
		{
			name: "binary",
			input: func() []byte {
				input := make([]byte, 256)
				for i := range input {
					input[i] = byte(i)
				}
				return input
			}(),
			expected: func() [256]uint64 {
				var expected [256]uint64
				for i := range expected {
					expected[i] = 1
				}
				return expected
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CalculateFrequencies(bytes.NewReader(test.input))
			if err != nil {
				t.Fatalf("CalculateFrequencies() error = %v", err)
			}
			if got != test.expected {
				t.Errorf("CalculateFrequencies() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestCalculateFrequenciesReadWithEOF(t *testing.T) {
	input := []byte("final bytes")
	got, err := CalculateFrequencies(&eofWithDataReader{data: input})
	if err != nil {
		t.Fatalf("CalculateFrequencies() error = %v", err)
	}

	var expected [256]uint64
	for _, b := range input {
		expected[b]++
	}
	if got != expected {
		t.Errorf("CalculateFrequencies() = %v, want %v", got, expected)
	}
}

func TestCalculateFrequenciesReaderFailure(t *testing.T) {
	readErr := errors.New("read failed")
	got, err := CalculateFrequencies(errorReader{err: readErr})
	if !errors.Is(err, readErr) {
		t.Fatalf("CalculateFrequencies() error = %v, want %v", err, readErr)
	}
	if got != ([256]uint64{}) {
		t.Errorf("CalculateFrequencies() frequencies = %v, want all zeros", got)
	}
}

func TestCalculateFrequenciesNoProgress(t *testing.T) {
	_, err := CalculateFrequencies(zeroReader{})
	if !errors.Is(err, io.ErrNoProgress) {
		t.Fatalf("CalculateFrequencies() error = %v, want %v", err, io.ErrNoProgress)
	}
}

func TestCalculateFrequenciesPartialDataReaderFailure(t *testing.T) {
	readErr := errors.New("read failed after partial data")
	got, err := CalculateFrequencies(&partialErrorReader{err: readErr})
	if !errors.Is(err, readErr) {
		t.Fatalf("CalculateFrequencies() error = %v, want %v", err, readErr)
	}

	var expected [256]uint64
	expected['a'] = 1
	expected['b'] = 1
	expected['c'] = 1
	if got != expected {
		t.Errorf("CalculateFrequencies() frequencies = %v, want %v", got, expected)
	}
}

type eofWithDataReader struct {
	data   []byte
	offset int
}

func (r *eofWithDataReader) Read(p []byte) (int, error) {
	if r.offset == len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.offset:])
	r.offset += n
	if r.offset == len(r.data) {
		return n, io.EOF
	}
	return n, nil
}

type errorReader struct {
	err error
}

func (r errorReader) Read([]byte) (int, error) {
	return 0, r.err
}

type -zeroReader struct{}

func (zeroReader) Read([]byte) (int, error) {
	return 0, nil
}

type partialErrorReader struct {
	err  error
	read bool
}

func (r *partialErrorReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, r.err
	}
	r.read = true
	return copy(p, []byte("abc")), nil
}
