# -*- coding: utf-8 -*-

import asyncio
import importlib
import os
import re
import tempfile
import time
from pathlib import Path
from random import choice
from typing import List

GO_MODULE_BASE: Path = Path(__file__).parent.parent.joinpath("components")
BENCHMARK_RESULTS_BASE: Path = Path(__file__).parent.parent.joinpath("components").joinpath("benchmarks").joinpath("results")

base_common_dir = "benchmarking._internal.common.md_py_common.py_common"

command_handling_package = base_common_dir + ".command_handling"
logger_package = base_common_dir + ".logging"
cli_package = base_common_dir + ".cli_framework"
file_handler_package = base_common_dir + ".handlers"

# Import the module
logger_module = importlib.import_module(logger_package)
command_handling_module = importlib.import_module(command_handling_package)
cli_module = importlib.import_module(cli_package)
file_handler_module = importlib.import_module(file_handler_package)

HoornLogger = logger_module.HoornLogger
LogType = logger_module.LogType
DefaultHoornLogOutput = logger_module.DefaultHoornLogOutput
CLIInterface = cli_module.CommandLineInterface
FileHandler = file_handler_module.FileHandler

#TODO - CLEAN UP CODE

def _benchmark(command_handler, num_secs: int, profile: str):
	commands = [
		"tool",
		"pprof",
		"-http=:8081",
		f"--seconds {num_secs}",
		f"http://localhost:6060/debug/pprof/{profile}"
	]

	command_handler.execute_command_v2("go", commands, shell=True, hide_console=False, keep_open=False)

def benchmark_command(command_handler):
	profile_options = {
		"profile": "CPU profile. This shows where your program is spending its CPU time.",
		"heap": "Memory allocations. This shows how your program is using memory, including where allocations are happening and how big they are.",
		"block": "Goroutine blocking events. This shows where goroutines are blocking, such as waiting for locks or channels.",
		"goroutine": "Stack traces of all current goroutines. This shows what every goroutine in your program is currently doing.",
		"threadcreate": "Stack traces that led to the creation of new OS threads. This can help you understand why your program is creating a lot of threads.",
		"mutex": "Contention profiles for mutexes. This shows where mutexes are causing contention between goroutines.",
		"trace": "Provides a trace of execution events in your Go program. You can visualize this data with the `go tool trace` command."
	}

	option_num = -1
	for option, description in profile_options.items():
		option_num += 1
		print(f"{option_num}. {option}: {description}")

	choice = int(input("Enter the number of the profile option: "))

	# Validate the choice
	if choice < 0 or choice >= len(profile_options):
		print("Invalid choice. Please try again.")
		time.sleep(0.5)
		return benchmark_command(command_handler)

	profile_option = list(profile_options.keys())[choice]

	num_secs = int(input("Enter the number of seconds to run the benchmark: "))

	_benchmark(command_handler, num_secs, profile_option)

def _get_benchmark_choice(go_files) -> Path:
	for i, file in enumerate(go_files):
		print(f"{i}. {file.name}")

	choice = int(input("Enter the number of the benchmark file: "))

	if choice < 0 or choice >= len(go_files):
		print("Invalid choice. Please try again.")
		time.sleep(0.5)
		return _get_benchmark_choice(go_files)

	return go_files[choice]

def benchmark_command_2(command_handler, file_helper):
	async def run_benchmark():
		benchmark_path = GO_MODULE_BASE.joinpath("benchmarks")
		go_files = file_helper.get_children_paths(benchmark_path, extension=".go")

		if len(go_files) == 0:
			print("No benchmark files found in the benchmarks directory.")
			return

		choice = _get_benchmark_choice(go_files)
		output_name = input("Enter the output name for the benchmark results: ")

		benchmark_file = benchmark_path.joinpath(choice.name)

		desired_count = input("Enter the number of times to run the benchmark (leave empty for 10): ")
		desired_count = int(desired_count) if desired_count else 10

		commands = [
			"test",
			f"\"{benchmark_file.__str__()}\"",
			"-run=^$",
			f"-count={desired_count}",
			"-bench=.",
			"-benchmem",
			"-cpuprofile",
			f"\"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.cpu.prof').resolve()}\"",
			"-memprofile",
			f"\"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.mem.prof').resolve()}\"",
			f"-o=\"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.exe').resolve()}\"",
			f"> \"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.txt').resolve()}\""
		]

		commands_2 = [
			"tool",
			"pprof",
			f"\"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.exe').resolve()}\"",
			f"\"{BENCHMARK_RESULTS_BASE.joinpath(f'{output_name}.cpu.prof')}\""
		]

		os.chdir(GO_MODULE_BASE)

		await command_handler.execute_command_v2_async("go", commands, hide_console=False, keep_open=False)
		await command_handler.execute_command_v2_async("go", commands_2, hide_console=False, keep_open=False)

	asyncio.run(run_benchmark())

def benchstat_compare_command(command_handler, file_helper):
	async def run_benchstat_compare():
		benchmark_result_path = GO_MODULE_BASE.joinpath("benchmarks").joinpath("results")
		interpreted_result_path = GO_MODULE_BASE.joinpath("benchmarks").joinpath("interpretation")
		results: List[Path] = file_helper.get_children_paths(benchmark_result_path, extension=".txt")

		if len(results) == 0:
			print("No benchmark results found in the results directory.")
			return

		results.sort(key=lambda x: x.name)
		print("Available benchmark results:")
		for i, result in enumerate(results):
			print(f"{i}. {result.name}")

		choice_1 = int(input("Enter the number of the first benchmark result: "))
		choice_2 = int(input("Enter the number of the second benchmark result: "))

		if choice_1 < 0 or choice_1 >= len(results) or choice_2 < 0 or choice_2 >= len(results):
			print("Invalid choice. Please try again.")
			return benchstat_compare_command(command_handler, file_helper)

		result_1 = results[choice_1]
		result_2 = results[choice_2]

		with open(interpreted_result_path.joinpath(f"{result_1.stem}-{result_2.stem}.txt"), "w+") as tmpfile:
			temp_file_path = tmpfile.name

		command = [
	        f"O=\"{result_1.resolve()}\"",
	        f"N=\"{result_2.resolve()}\"",
			f"> \"{temp_file_path}\""
		]

		await command_handler.execute_command_v2_async("benchstat", command, hide_console=False, keep_open=False)

		print(f"Benchstat comparison (cleaned) written to: {temp_file_path}")

	asyncio.run(run_benchstat_compare())

def run_pprof_command(command_handler, file_helper):
	exes: List[Path] = file_helper.get_children_paths(BENCHMARK_RESULTS_BASE, extension=".exe")
	if len(exes) == 0:
		print("No benchmark executables found in the benchmark results directory.")
		return

	for i, exe in enumerate(exes):
		print(f"{i}. {exe.name}")

	chosen_executable = int(input("Enter the number of the benchmark executable: "))
	if chosen_executable < 0 or chosen_executable >= len(exes):
		print("Invalid choice. Please try again.")
		return run_pprof_command(command_handler, file_helper)

	executable = exes[chosen_executable]
	cpu_result = BENCHMARK_RESULTS_BASE.joinpath(f"{executable.stem}.cpu.prof")

	commands = [
        "tool",
        "pprof",
        f"\"{executable.resolve()}\"",
	    f"\"{cpu_result.resolve()}\""
	]

	command_handler.execute_command_v2("go", commands, shell=True, hide_console=False, keep_open=True)

def run_pprof_web_command(command_handler, file_helper):
	exes: List[Path] = file_helper.get_children_paths(BENCHMARK_RESULTS_BASE, extension=".exe")
	if len(exes) == 0:
		print("No benchmark executables found in the benchmark results directory.")
		return

	for i, exe in enumerate(exes):
		print(f"{i}. {exe.name}")

	chosen_executable = int(input("Enter the number of the benchmark executable: "))
	if chosen_executable < 0 or chosen_executable >= len(exes):
		print("Invalid choice. Please try again.")
		return run_pprof_command(command_handler, file_helper)

	executable = exes[chosen_executable]
	cpu_result = BENCHMARK_RESULTS_BASE.joinpath(f"{executable.stem}.cpu.prof")

	commands = [
		"tool",
		"pprof",
		"-http=\":8080\"",
		f"\"{executable.resolve()}\"",
		f"\"{cpu_result.resolve()}\""
	]

	command_handler.execute_command_v2("go", commands, hide_console=False, keep_open=True, shell=True)


if __name__ == "__main__":
	logger = HoornLogger(min_level=LogType.DEBUG, outputs=[DefaultHoornLogOutput()])
	command_handler = command_handling_module.CommandHelper(logger)
	cli_interface = CLIInterface(logger)
	file_handler = FileHandler()

	cli_interface.add_command(["benchmark", "bm"], action=benchmark_command, description="Starts the benchmarking.", arguments=[command_handler])
	cli_interface.add_command(["benchmark-2", "bm-2"], action=benchmark_command_2, description="Starts the benchmarking suite.", arguments=[command_handler, file_handler])
	cli_interface.add_command(["benchstat-compare", "bsc"], action=benchstat_compare_command, description="Compares two benchmark results.", arguments=[command_handler, file_handler])
	cli_interface.add_command(["pprof", "pp"], action=run_pprof_command, description="Starts the pprof tool.", arguments=[command_handler, file_handler])
	cli_interface.add_command(["pprof-web", "pp-web"], action=run_pprof_web_command, description="Starts the pprof tool with web interface.", arguments=[command_handler, file_handler])

	cli_interface.start_listen_loop()


