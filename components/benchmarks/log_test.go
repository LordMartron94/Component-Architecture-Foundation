package benchmarks

import (
	"os"
	"sync"
	"testing"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/benchmarks/_internal"
)

func BenchmarkLoggingSuite(b *testing.B) {
	b.Run("Log Debug 1x", func(b *testing.B) {
		benchmarkLogDebug(b, 1)
	})

	b.Run("Log Debug 100x", func(b *testing.B) {
		benchmarkLogDebug(b, 100)
	})

	b.Run("Log Debug 10000x", func(b *testing.B) {
		benchmarkLogDebug(b, 10000)
	})
}

func benchmarkLogDebug(b *testing.B, iterations int) {
	tmpFile, err := os.CreateTemp("", "benchmark-log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	originalStdout := os.Stdout

	os.Stdout = tmpFile

	shutdownSignal := make(chan struct{})
	wg := &sync.WaitGroup{}

	logger := _internal.GetLogger("Component-Based Architecture Foundation - Benchmarks", shutdownSignal, wg)

	b.ResetTimer() // Start timer after logger initialization
	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			logger.Debug("Benchmark debug message", false, "")
		}
	}

	close(shutdownSignal)
	wg.Wait()

	os.Stdout = originalStdout
}
