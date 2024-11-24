package output

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/formatting"
)

type FileHoornLogOutput struct {
	logDirectory    string
	maxLogsToKeep   int
	createDirectory bool
	useCombined     bool

	logsToWrite []*common.HoornLog

	validatedDirectories map[string]bool
}

func NewFileHoornLogOutput(logDirectory string, maxLogsToKeep int, useCombined bool) *FileHoornLogOutput {
	var fileHoornLogOutput = &FileHoornLogOutput{
		logDirectory:         filepath.Clean(logDirectory),
		maxLogsToKeep:        maxLogsToKeep,
		createDirectory:      true,
		useCombined:          useCombined,
		logsToWrite:          make([]*common.HoornLog, 0, 300),
		validatedDirectories: make(map[string]bool, 30),
	}

	fileHoornLogOutput.initialize()

	return fileHoornLogOutput
}

func NewFileHoornLogOutputWithoutCreateDir(logDirectory string, maxLogsToKeep int, useCombined bool) *FileHoornLogOutput {
	var fileHoornLogOutput = &FileHoornLogOutput{
		logDirectory:         filepath.Clean(logDirectory),
		maxLogsToKeep:        maxLogsToKeep,
		createDirectory:      false,
		useCombined:          useCombined,
		logsToWrite:          make([]*common.HoornLog, 0, 300),
		validatedDirectories: make(map[string]bool, 30),
	}

	fileHoornLogOutput.initialize()

	return fileHoornLogOutput
}

func (fhl *FileHoornLogOutput) initialize() {
	fhl.validateDirectory(fhl.logDirectory)
	fhl.incrementLogs()
}

func (fhl *FileHoornLogOutput) validateDirectory(directory string) error {
	if fhl.validatedDirectories[directory] {
		return nil
	}

	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		if fhl.createDirectory {
			errDir := os.MkdirAll(directory, 0755)
			if errDir != nil {
				return errDir
			}

			fhl.validatedDirectories[directory] = true
			return nil
		}
		return fmt.Errorf("log directory %v does not exist", directory)
	}

	fhl.validatedDirectories[directory] = true
	return nil
}

// getSubDirectories will retrieve a list of subdirectories
func (fhl *FileHoornLogOutput) getSubDirectories() ([]os.DirEntry, error) {
	return os.ReadDir(fhl.logDirectory)
}

// handleLogFile will handle each log file
func handleLogFile(dirPath, file string, maxLogsToKeep int) error {
	extension := filepath.Ext(file)
	name := strings.TrimSuffix(file, extension)
	splitName := strings.Split(name, "_")
	logNumber, err := strconv.Atoi(splitName[len(splitName)-1])
	if err != nil {
		return err
	}

	if logNumber+1 > maxLogsToKeep {
		err := os.Remove(file)
		if err != nil {
			return err
		}
		return nil
	}

	err = os.Rename(file, filepath.Join(dirPath, fmt.Sprintf("log_%v.txt", logNumber+1)))
	if err != nil {
		return err
	}

	return nil
}

func (fhl *FileHoornLogOutput) handleDir(dir string) error {
	children, err := getFileChildrenPaths(dir, ".txt")
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(children)))

	for _, file := range children {
		err := handleLogFile(dir, file, fhl.maxLogsToKeep)
		if err != nil {
			return err
		}
	}

	return nil
}

// incrementLogs will call smaller functions to increment logs
func (fhl *FileHoornLogOutput) incrementLogs() error {
	subDirectories, err := fhl.getSubDirectories()
	if err != nil {
		return err
	}

	if fhl.useCombined {
		err := fhl.handleDir(fhl.logDirectory)

		if err != nil {
			return err
		}
	}

	for _, subDir := range subDirectories {
		if subDir.IsDir() {
			err := fhl.handleDir(filepath.Join(fhl.logDirectory, subDir.Name()))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fhl *FileHoornLogOutput) getPathToLogTo(logSeparator []byte) string {
	directory := fhl.logDirectory

	if len(logSeparator) > 0 {
		directory = filepath.Join(fhl.logDirectory, string(logSeparator))
	}

	fhl.validateDirectory(directory)

	return filepath.Join(directory, fmt.Sprintf("log_%v.txt", 1))
}

func getFileChildrenPaths(directory string, extension string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(directory)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() && filepath.Ext(file.Name()) == extension {
			files = append(files, filepath.Join(directory, file.Name()))
		}
	}
	return files, nil
}

func (fhl *FileHoornLogOutput) writeLogs(separators [][]byte, messages [][][]byte) {
	for i, separator := range separators {
		var logDirectory = fhl.getPathToLogTo(separator)
		logsAssociatedWithSeparator := messages[i]

		for _, msg := range logsAssociatedWithSeparator {
			f, err := os.OpenFile(logDirectory, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}

			formatted := append(msg, []byte("\n")...)

			if _, err := f.Write(formatted); err != nil {
				log.Fatal(err)
			}

			f.Close()
		}
	}
}

func (fhl *FileHoornLogOutput) Output(hoornLog *common.HoornLog) {
	fhl.logsToWrite = append(fhl.logsToWrite, hoornLog)
}

func alternativeContains(container [][]byte, search []byte) bool {
	for _, element := range container {
		if bytes.Equal(element, search) {
			return true
		}
	}

	return false
}

func (fhl *FileHoornLogOutput) Save() {
	textFormatter := formatting.NewHoornLogTextFormatter()
	combinedTextFormatter := formatting.NewHoornLogCombinedTextFormatter(*textFormatter)

	separators := make([][]byte, 0, 20)
	messages := make([][][]byte, 0, 20)

	separators = append(separators, []byte(""))
	indexOfCombinedSeparator := 0

	for _, hoornLog := range fhl.logsToWrite {
		formattedLog := textFormatter.Format(hoornLog)

		if !alternativeContains(separators, hoornLog.LogSeparator) {
			separators = append(separators, hoornLog.LogSeparator)
		}

		indexOfSeparator := sort.Search(len(separators), func(i int) bool { return bytes.Equal(separators[i], hoornLog.LogSeparator) })

		if indexOfSeparator >= len(messages) {
			newMessages := make([][][]byte, indexOfSeparator+1)
			copy(newMessages, messages)
			messages = newMessages
		}

		if messages[indexOfSeparator] == nil {
			messages[indexOfSeparator] = make([][]byte, 0, 10)
		}

		messages[indexOfSeparator] = append(messages[indexOfSeparator], formattedLog)

		if fhl.useCombined {
			formattedLog = combinedTextFormatter.Format(hoornLog)
			messages[indexOfCombinedSeparator] = append(messages[indexOfCombinedSeparator], formattedLog)
		}
	}

	fhl.writeLogs(separators, messages)

	fhl.logsToWrite = make([]*common.HoornLog, 0, 300)
}
