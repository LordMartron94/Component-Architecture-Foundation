package scanning

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/shared"
)

// Scanner is a custom scanner that reads from an io.Reader and splits the input
// based on a custom delimiter string.
type Scanner struct {
	Logger    *logging.HoornLogger
	reader    *bufio.Reader
	delimiter []byte
	buffer    []byte
	eof       bool
}

// NewScanner creates a new Scanner with the given io.Reader and delimiter string.
func NewScanner(r io.Reader, delimiter string, logger *logging.HoornLogger) *Scanner {
	return &Scanner{
		reader:    bufio.NewReader(r),
		delimiter: []byte(delimiter),
		buffer:    make([]byte, 0),
		eof:       false,
		Logger:    logger,
	}
}

// Scan advances the Scanner to the next token, returning true if a token was found,
// and false otherwise (either because the end of the input was reached or an error occurred).
func (s *Scanner) Scan(shutdownChan <-chan struct{}) (bool, error) {
	s.buffer = s.buffer[:0]
	//s.Logger.Debug("Starting to scan for next token", false, shared.ScannerComponentName)
	//defer s.Logger.Debug("Finished scanning for next token", false, shared.ScannerComponentName)

	for {
		//s.Logger.Debug("Scanning for next token", false, shared.ScannerComponentName)

		// Select statement to handle shutdown signal and check connection status
		select {
		case <-shutdownChan:
			// Shutdown signal received, return immediately
			s.Logger.Info("Scanner shutting down", false, shared.ScannerComponentName)
			return false, nil
		default:
			// Read byte by byte until the delimiter or EOF is reached
			for {
				select {
				case <-shutdownChan:
					// Shutdown signal received, return immediately
					s.Logger.Info("Scanner shutting down", false, shared.ScannerComponentName)
					return false, nil
				default:
					b, err := s.reader.ReadByte()
					if err != nil {
						if err == io.EOF {
							s.eof = true
							// Return true if there's any remaining data in the buffer
							s.Logger.Info("Returning remaining data in buffer, end of file", false, shared.ScannerComponentName)
							return len(s.buffer) > 0, err
						}

						s.Logger.Warn(fmt.Sprintf("Error reading bytes: '%s", err), false, shared.ScannerComponentName)
						return false, err
					}

					s.buffer = append(s.buffer, b)

					if bytes.HasSuffix(s.buffer, s.delimiter) {
						// Found the delimiter, trim it and return
						s.buffer = bytes.TrimSuffix(s.buffer, s.delimiter)
						return len(s.buffer) > 0 || s.eof, nil
					}
				}
			}
		}
	}
}

// Bytes returns the current token as bytes.
func (s *Scanner) Bytes() []byte {
	//s.Logger.Debug(fmt.Sprintf("Current bytes: %s", s.buffer), false, shared.ScannerComponentName)
	return append([]byte(nil), s.buffer...) // Create a copy to avoid data races
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
