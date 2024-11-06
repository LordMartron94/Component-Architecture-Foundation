package scanning

import (
	"bufio"
	"bytes"
	"io"
)

// Scanner is a custom scanner that reads from an io.Reader and splits the input
// based on a custom delimiter string.
type Scanner struct {
	reader    *bufio.Reader
	delimiter []byte
	buffer    []byte
	eof       bool
}

// NewScanner creates a new Scanner with the given io.Reader and delimiter string.
func NewScanner(r io.Reader, delimiter string) *Scanner {
	return &Scanner{
		reader:    bufio.NewReader(r),
		delimiter: []byte(delimiter),
		buffer:    make([]byte, 0),
		eof:       false,
	}
}

// Scan advances the Scanner to the next token, returning true if a token was found,
// and false otherwise (either because the end of the input was reached or an error occurred).
func (s *Scanner) Scan() (bool, error) {
	for {
		// Read byte by byte until the delimiter or EOF is reached
		for {
			b, err := s.reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					s.eof = true
					return len(s.buffer) > 0, err // Return true if there's any remaining data
				}

				return false, err
			}

			s.buffer = append(s.buffer, b)

			if bytes.HasSuffix(s.buffer, s.delimiter) {
				break // Found the delimiter
			}
		}

		// Trim the delimiter from the buffer
		s.buffer = bytes.TrimSuffix(s.buffer, s.delimiter)

		// Return true even for empty tokens if it's the last one
		return len(s.buffer) > 0 || s.eof, nil
	}
}

// Bytes returns the current token as bytes.
func (s *Scanner) Bytes() []byte {
	return s.buffer
}

// Text returns the current token as a string.
func (s *Scanner) Text() string {
	return string(s.buffer)
}

// Reset resets the Scanner to the beginning of the input. This requires the underlying reader to be seekable.
func (s *Scanner) Reset(r io.Reader) {
	if seeker, ok := r.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
		s.reader.Reset(r)
		s.buffer = s.buffer[:0]
		s.eof = false
	} else {
		panic("Reader is not seekable") // or handle gracefully
	}
}
