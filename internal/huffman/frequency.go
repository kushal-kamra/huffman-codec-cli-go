package huffman

import "io"

const frequencyBufferSize = 32 * 1024

func CalculateFrequencies(input io.Reader) ([256]uint64, error) {
	var frequencies [256]uint64
	buffer := make([]byte, frequencyBufferSize)
	emptyReads := 0

	for {
		n, err := input.Read(buffer)
		for _, b := range buffer[:n] {
			frequencies[b]++
		}

		if err != nil {
			if err == io.EOF {
				return frequencies, nil
			}
			return frequencies, err
		}

		if n == 0 {
			emptyReads++
			if emptyReads >= 100 {
				return frequencies, io.ErrNoProgress
			}
			continue
		}
		emptyReads = 0
	}
}
