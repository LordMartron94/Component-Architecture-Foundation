package benchmarks

import (
	"os"
	"testing"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/benchmarks/_internal"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
)

func BenchmarkLoggingSuite(b *testing.B) {
	logger := _internal.GetLogger("Component-Based Architecture Foundation - Benchmarks")

	b.Run("LogDebug", func(b *testing.B) {
		benchmarkLogDebug(b, &logger)
	})
}

func benchmarkLogDebug(b *testing.B, logger *logging.HoornLogger) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "benchmark-log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file

	// Store original output
	originalStdout := os.Stdout

	// Redirect console output to the temporary file
	os.Stdout = tmpFile

	for i := 0; i < b.N; i++ {
		logger.Debug("Benchmark log debug message", false, "")
	}

	// Restore original output
	os.Stdout = originalStdout
}
