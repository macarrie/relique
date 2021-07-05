package rsync

import (
	"bufio"
	"os"
	"regexp"
	"strconv"

	log "github.com/macarrie/relique/internal/logging"

	"github.com/pkg/errors"
)

var regexProgress = regexp.MustCompile(`.*\s*(?P<speed>[\d,.]+)kB/s.*\(xfr#\d+, to-chk=(?P<remaining>\d+)/(?P<total>\d+)\)`)

// Progress contains rsync job progress
type Progress struct {
	// Total number of files to sync
	Total int
	// Number of remaining files
	Remaining int
	// Files already handled
	Current int
	// Sync progress of all files in percent
	Percent float32
	// Current transfer speed
	Speed float64
}

func (p *Progress) GetFromRsyncLog(path string) error {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Getting rsync task progress from log file")

	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "cannot read log file")
	}
	defer file.Close()

	var matchingLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		line := string(lineBytes)

		if regexProgress.Match(lineBytes) {
			matchingLines = append(matchingLines, line)
		}
	}

	if len(matchingLines) == 0 {
		return nil
	}

	// Get only last line result, no need to parse previous progress lines
	progressMap := getMatchMap(regexProgress, matchingLines[len(matchingLines)-1])
	if remaining, ok := progressMap["remaining"]; ok {
		p.Remaining, _ = strconv.Atoi(remaining)
	}
	if total, ok := progressMap["total"]; ok {
		p.Total, _ = strconv.Atoi(total)
	}
	if speed, ok := progressMap["speed"]; ok {
		p.Speed, _ = strconv.ParseFloat(stripCommas(speed), 64)
	}

	if p.Total != 0 {
		p.Current = p.Total - p.Remaining
		p.Percent = (float32(p.Current) / float32(p.Total)) * 100
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "error during log file scan")
	}

	return nil
}
