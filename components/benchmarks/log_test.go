package benchmarks

import (
	"os"
	"testing"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/benchmarks/_internal"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
)

func BenchmarkLoggingSuite(b *testing.B) {
	logger := _internal.GetLogger("Component-Based Architecture Foundation - Benchmarks")

	b.Run("Log Debug 1x", func(b *testing.B) {
		benchmarkLogDebug1(b, &logger)
	})

	b.Run("Log Debug 100x", func(b *testing.B) {
		benchmarkLogDebug100(b, &logger)
	})

	b.Run("Log Debug 10000x", func(b *testing.B) {
		benchmarkLogDebug10k(b, &logger)
	})
}

func benchmarkLogDebug1(b *testing.B, logger *logging.HoornLogger) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "benchmark-log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	originalStdout := os.Stdout

	os.Stdout = tmpFile

	for i := 0; i < b.N; i++ {
		sendDebug(logger)
	}

	os.Stdout = originalStdout
}

func benchmarkLogDebug100(b *testing.B, logger *logging.HoornLogger) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "benchmark-log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	originalStdout := os.Stdout

	os.Stdout = tmpFile

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			sendDebug(logger)
		}
	}

	os.Stdout = originalStdout
}

func benchmarkLogDebug10k(b *testing.B, logger *logging.HoornLogger) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "benchmark-log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	originalStdout := os.Stdout

	os.Stdout = tmpFile

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			sendDebug(logger)
		}
	}

	os.Stdout = originalStdout
}

func sendDebug(logger *logging.HoornLogger) {
	logger.Debug("Benchmark log debug message", false, "")
}
