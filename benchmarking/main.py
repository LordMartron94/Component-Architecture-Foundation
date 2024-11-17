import importlib
import time

base = "benchmarking._internal.common.md_py_common.py_common"

command_handling_package = base + ".command_handling"
logger_package = base + ".logging"
cli_package = base + ".cli_framework"

# Import the module
logger_module = importlib.import_module(logger_package)
command_handling_module = importlib.import_module(command_handling_package)
cli_module = importlib.import_module(cli_package)

HoornLogger = logger_module.HoornLogger
LogType = logger_module.LogType
DefaultHoornLogOutput = logger_module.DefaultHoornLogOutput
CLIInterface = cli_module.CommandLineInterface

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


if __name__ == "__main__":
	logger = HoornLogger(min_level=LogType.DEBUG, outputs=[DefaultHoornLogOutput()])
	command_handler = command_handling_module.CommandHelper(logger)
	cli_interface = CLIInterface(logger)

	cli_interface.add_command(["benchmark", "bm"], action=benchmark_command, description="Starts the benchmarking.", arguments=[command_handler])

	cli_interface.start_listen_loop()


