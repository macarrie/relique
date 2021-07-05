package rsync

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/pkg/errors"
)

var regexNumberOfFiles = regexp.MustCompile(`Number of files: (?P<nb>\d+) \(reg: (?P<reg>\d+), dir: (?P<dir>\d+)\)`)
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
	NumberOfFiles int
	// Number of regular files handled (excluding directories, symlinks, etc...)
	NumberOfRegularFiles int
	// Number of directories handled (excluding directories, symlinks, etc...)
	NumberOfDirectories int
	// Number of deleted files on destination
	NumberOfDeletedFiles int
	// Number of created files on destination
	NumberOfCreatedFiles int
	// Number of created files on destination (excluding folders and symlinks, etc...)
	NumberOfCreatedRegularFiles int
	// Total file size is the total sum of all file sizes in the transfer.  This  does  not  count any size for directories or special files, but does include the size of symlinks.
	TotalFileSize int64
	// Total  transferred  file  size is the total sum of all files sizes for just the transferred files.
	TotalTransferredFileSize int64
	// Literal data is how much unmatched file-update data we had to send to the receiver  for  it to recreate the updated files.
	LiteralData int64
	// Matched data is how much data the receiver got locally when recreating the updated files.
	MatchedData int64
	// File list size (in bytes)
	FileListSize int64
	// Time to generate file list (in seconds)
	FileListGenerationTime float64
	// Time to transfer file list (in seconds)
	FileListTransferTime float64
	// Total  bytes sent is the count of all the bytes that rsync sent from the client side to the server side.
	TotalBytesSent int64
	// Total bytes received is the count of all non-message  bytes  that  rsync  received  by  the client  side from the server side.  "Non-message" bytes means that we don’t count the bytes for a verbose message that the server sent to us, which makes the stats more consistent.
	TotalBytesReceived int64
	// Transfer speed rate (in bytes/second)
	TransferSpeed float64
	// Transfer speedup thanks to rsync diff algorithm
	TransferSpeedup float64
}

func (s *Stats) GetFromRsyncLog(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Getting rsync task stats from log file")
	logContentsBuffer, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "cannot read log file")
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
