package importer

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"parser/data"
	"slices"
	"strings"
)

type Importer struct {
	path       *string
	skipHeader *bool
	fieldIndex *int
	domainData []data.CustomerDomainData
}

// NewImporter returns a new Importer that reads from file at specified path.
func NewImporter(filePath *string, skipHeader *bool, fieldIndex *int) (*Importer, error) {
	if _, err := os.Stat(*filePath); err != nil {
		return nil, fmt.Errorf("file does exists or cannot be read %s: %w", *filePath, err)
	}

	return &Importer{
		path:       filePath,
		skipHeader: skipHeader,
		fieldIndex: fieldIndex,
	}, nil
}

// ImportDomainData reads and returns sorted customer domain data from CSV file.
func (ci Importer) ImportDomainData() ([]data.CustomerDomainData, error) {
	file, err := os.Open(*ci.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", *ci.path, err)
	}

	slog.Info("opened file", "path", *ci.path)

	defer func() {
		if err := file.Close(); err != nil {
			slog.Warn("failed to close file", "path", *ci.path, "err", err)
		} else {
			slog.Info("closed file", "path", *ci.path)
		}
	}()

	csvReader := csv.NewReader(file)
	domainMap := make(map[string]uint64)

	// skip first line with headers
	if *ci.skipHeader {
		slog.Info("skipping first line")
		if _, err := csvReader.Read(); err != nil && err != io.EOF {
			slog.Warn("failed to read line of csv", "err", err)
		}
	}

	// read each line and get domain by parsing email
	for line, err := csvReader.Read(); err != io.EOF; line, err = csvReader.Read() {
		if err != nil {
			slog.Warn("failed to read line of csv", "err", err)
		}

		if len(line) <= int(*ci.fieldIndex) {
			slog.Warn("index of field doesnt exists", "line", line)
			continue
		}

		if domain := parseEmail(line[*ci.fieldIndex]); domain != "" {
			domainMap[domain] += 1
		}
	}

	domainData := mapToSortedDomainData(domainMap)

	return domainData, nil
}

func parseEmail(s string) string {
	email, domain, isFound := strings.Cut(s, "@")

	if email == "" || !isFound || !strings.Contains(domain, ".") {
		slog.Warn("invalid email address", "email", email)
	}

	return domain
}

func mapToSortedDomainData(m map[string]uint64) []data.CustomerDomainData {
	domainData := make([]data.CustomerDomainData, 0, len(m))
	for domain, count := range m {
		domainData = append(domainData, data.CustomerDomainData{
			Domain:           domain,
			CustomerQuantity: count,
		})
	}

	slices.SortFunc(domainData, func(l, r data.CustomerDomainData) int {
		return cmp.Compare(l.Domain, r.Domain)
	})

	return domainData
}
