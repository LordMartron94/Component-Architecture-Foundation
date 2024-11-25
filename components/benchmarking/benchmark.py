from pathlib import Path

from components.benchmarking.tool.cli_bench.api import start_golang_cli

if __name__ == "__main__":
	module_root: Path = Path(__file__).parent.parent
	benchmarks_path: Path = module_root.joinpath("benchmarks")
	benchmark_results_path: Path = benchmarks_path.joinpath("results")
	benchmark_interpretation_path: Path = benchmarks_path.joinpath("interpretation")

	start_golang_cli(
		str(module_root),
		str(benchmarks_path),
        str(benchmark_results_path),
        str(benchmark_interpretation_path)
	)

