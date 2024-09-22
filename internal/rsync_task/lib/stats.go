package rsync_lib

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"
)

var regexNumberOfFiles = regexp.MustCompile(`Number of files: (?P<nb>\d+) \(reg: (?P<reg>\d+)(, dir: (?P<dir>\d+))?\)`)
var regexNumberOfCreatedFiles = regexp.MustCompile(`Number of created files: (?P<nb>\d+) \(reg: (?P<reg>\d+)\)`)
var regexNumberOfDeletedFiles = regexp.MustCompile(`Number of deleted files: (?P<nb>\d+)`)
var regexTotalFileSize = regexp.MustCompile(`Total file size: (?P<size>[\d,]+) bytes`)
var regexTotalTransferredFileSize = regexp.MustCompile(`Total transferred file size: (?P<size>[\d,]+) bytes`)
var regexLiteralData = regexp.MustCompile(`Literal data: (?P<size>[\d,]+) bytes`)
var regexMatchedData = regexp.MustCompile(`Matched data: (?P<size>[\d,]+) bytes`)
var regexFileListSize = regexp.MustCompile(`File list size: (?P<size>[\d,]+)`)
var regexFileListGenerationTime = regexp.MustCompile(`File list generation time: (?P<nb>[\d.]+) seconds`)
var regexFileListTransferTime = regexp.MustCompile(`File list transfer time: (?P<nb>[\d.]+) seconds`)
var regexTotalBytesSent = regexp.MustCompile(`Total bytes sent: (?P<nb>[\d,]+)`)
var regexTotalBytesReceived = regexp.MustCompile(`Total bytes received: (?P<nb>[\d,]+)`)
var regexTransferSpeed = regexp.MustCompile(`sent .* bytes\s*received .* bytes\s*(?P<nb>[\d,.]+) bytes/sec`)
var regexSpeedup = regexp.MustCompile(`total size is .*\s*speedup is\s*(?P<nb>[\d,\.]+)`)

// Stats contains rsync job statistics
type Stats struct {
	// Number of files is the count  of  all  "files"  (in  the  generic  sense),  which  includes directories, symlinks, etc.
	NumberOfFiles int `toml:"number_of_files"`
	// Number of regular files handled (excluding directories, symlinks, etc...)
	NumberOfRegularFiles int `toml:"number_of_regular_files"`
	// Number of directories handled (excluding directories, symlinks, etc...)
	NumberOfDirectories int `toml:"number_of_directories"`
	// Number of deleted files on destination
	NumberOfDeletedFiles int `toml:"number_of_deleted_files"`
	// Number of created files on destination
	NumberOfCreatedFiles int `toml:"number_of_created_files"`
	// Number of created files on destination (excluding folders and symlinks, etc...)
	NumberOfCreatedRegularFiles int `toml:"number_of_created_regular_files"`
	// Total file size is the total sum of all file sizes in the transfer.  This  does  not  count any size for directories or special files, but does include the size of symlinks.
	TotalFileSize int64 `toml:"total_file_size"`
	// Total  transferred  file  size is the total sum of all files sizes for just the transferred files.
	TotalTransferredFileSize int64 `toml:"total_transferred_file_size"`
	// Literal data is how much unmatched file-update data we had to send to the receiver  for  it to recreate the updated files.
	LiteralData int64 `toml:"literal_data"`
	// Matched data is how much data the receiver got locally when recreating the updated files.
	MatchedData int64 `toml:"matched_data"`
	// File list size (in bytes)
	FileListSize int64 `toml:"file_list_size"`
	// Time to generate file list (in seconds)
	FileListGenerationTime float64 `toml:"file_list_generation_time"`
	// Time to transfer file list (in seconds)
	FileListTransferTime float64 `toml:"file_list_transfer_time"`
	// Total  bytes sent is the count of all the bytes that rsync sent from the client side to the server side.
	TotalBytesSent int64 `toml:"total_bytes_sent"`
	// Total bytes received is the count of all non-message  bytes  that  rsync  received  by  the client  side from the server side.  "Non-message" bytes means that we don’t count the bytes for a verbose message that the server sent to us, which makes the stats more consistent.
	TotalBytesReceived int64 `toml:"total_bytes_received"`
	// Transfer speed rate (in bytes/second)
	TransferSpeed float64 `toml:"-"`
	// Transfer speedup thanks to rsync diff algorithm
	TransferSpeedup float64 `toml:"-"`
}

func (s *Stats) GetFromRsyncLog(path string) error {
	slog.With(
		slog.String("path", path),
	).Debug("Getting rsync task stats from log file")

	logContentsBuffer, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read log file: %w", err)
	}

	logContents := string(logContentsBuffer)

	nbFiles := getMatchMap(regexNumberOfFiles, logContents)
	if nbOfFiles, ok := nbFiles["nb"]; ok {
		s.NumberOfFiles, _ = strconv.Atoi(nbOfFiles)
	}
	if nbOfRegularFiles, ok := nbFiles["reg"]; ok {
		s.NumberOfRegularFiles, _ = strconv.Atoi(nbOfRegularFiles)
	}
	if nbOfDirectories, ok := nbFiles["dir"]; ok {
		s.NumberOfDirectories, _ = strconv.Atoi(nbOfDirectories)
	}

	nbCreatedFiles := getMatchMap(regexNumberOfCreatedFiles, logContents)
	if nbOfCreatedFiles, ok := nbCreatedFiles["nb"]; ok {
		s.NumberOfCreatedFiles, _ = strconv.Atoi(nbOfCreatedFiles)
	}
	if nbOfCreatedRegularFiles, ok := nbCreatedFiles["reg"]; ok {
		s.NumberOfCreatedRegularFiles, _ = strconv.Atoi(nbOfCreatedRegularFiles)
	}

	nbDeletedFiles := getMatchMap(regexNumberOfDeletedFiles, logContents)
	if nbOfDeletedFiles, ok := nbDeletedFiles["nb"]; ok {
		s.NumberOfDeletedFiles, _ = strconv.Atoi(nbOfDeletedFiles)
	}

	totalFileSizeResults := getMatchMap(regexTotalFileSize, stripCommas(logContents))
	if totalFileSize, ok := totalFileSizeResults["size"]; ok {
		s.TotalFileSize, _ = strconv.ParseInt(totalFileSize, 10, 64)
	}

	totalTransferredFileSizeResults := getMatchMap(regexTotalTransferredFileSize, stripCommas(logContents))
	if totalTransferredSize, ok := totalTransferredFileSizeResults["size"]; ok {
		s.TotalTransferredFileSize, _ = strconv.ParseInt(totalTransferredSize, 10, 64)
	}

	literalDataResults := getMatchMap(regexLiteralData, stripCommas(logContents))
	if literalData, ok := literalDataResults["size"]; ok {
		s.LiteralData, _ = strconv.ParseInt(literalData, 10, 64)
	}

	matchedDataResults := getMatchMap(regexMatchedData, stripCommas(logContents))
	if matchedData, ok := matchedDataResults["size"]; ok {
		s.MatchedData, _ = strconv.ParseInt(matchedData, 10, 64)
	}

	fileListSizeResults := getMatchMap(regexFileListSize, stripCommas(logContents))
	if fileListSize, ok := fileListSizeResults["size"]; ok {
		s.FileListSize, _ = strconv.ParseInt(fileListSize, 10, 64)
	}

	fileListGenTimeResults := getMatchMap(regexFileListGenerationTime, stripCommas(logContents))
	if fileListGenTime, ok := fileListGenTimeResults["nb"]; ok {
		s.FileListGenerationTime, _ = strconv.ParseFloat(fileListGenTime, 64)
	}

	fileListTransferTimeResults := getMatchMap(regexFileListTransferTime, stripCommas(logContents))
	if fileListTransferTime, ok := fileListTransferTimeResults["nb"]; ok {
		s.FileListTransferTime, _ = strconv.ParseFloat(fileListTransferTime, 64)
	}

	totalBytesSentResult := getMatchMap(regexTotalBytesSent, stripCommas(logContents))
	if totalBytesSent, ok := totalBytesSentResult["nb"]; ok {
		s.TotalBytesSent, _ = strconv.ParseInt(totalBytesSent, 10, 64)
	}

	totalBytesReceivedResult := getMatchMap(regexTotalBytesReceived, stripCommas(logContents))
	if totalBytesReceived, ok := totalBytesReceivedResult["nb"]; ok {
		s.TotalBytesReceived, _ = strconv.ParseInt(totalBytesReceived, 10, 64)
	}

	transferSpeedResults := getMatchMap(regexTransferSpeed, stripCommas(logContents))
	if transferSpeed, ok := transferSpeedResults["nb"]; ok {
		s.TransferSpeed, _ = strconv.ParseFloat(transferSpeed, 64)
	}

	speedUpResults := getMatchMap(regexSpeedup, stripCommas(logContents))
	if speedup, ok := speedUpResults["nb"]; ok {
		s.TransferSpeedup, _ = strconv.ParseFloat(speedup, 64)
	}

	return nil
}

func getMatchMap(r *regexp.Regexp, s string) map[string]string {
	matches := r.FindStringSubmatch(s)
	m := make(map[string]string)

	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(matches) {
			m[name] = matches[i]
		}
	}

	return m
}

func stripCommas(str string) string {
	return strings.ReplaceAll(str, ",", "")
}

func LoadStatsFromFile(file string) (s Stats, err error) {
	slog.Debug("Loading rsync stats from file", slog.String("path", file))

	f, err := os.Open(file)
	if err != nil {
		return Stats{}, fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("cannot close file correctly: %w", err)
		}
	}()

	content, _ := io.ReadAll(f)

	var st Stats
	if err := toml.Unmarshal(content, &st); err != nil {
		return Stats{}, fmt.Errorf("cannot parse toml file: %w", err)
	}

	return st, nil
}
