package importer

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestImportData(t *testing.T) {
	path := "./test_data.csv"
	skipHeader := true
	fieldIndex := 2
	importer, err := NewImporter(&path, &skipHeader, &fieldIndex)
	if err != nil {
		t.Error(err)
	}

	_, err = importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}
}

func TestImportDataSort(t *testing.T) {
	sortedDomains := []string{"360.cn", "acquirethisname.com", "blogtalkradio.com", "chicagotribune.com", "cnet.com", "cyberchimps.com", "github.io", "hubpages.com", "rediff.com", "statcounter.com"}
	path := "./test_data.csv"
	skipHeader := true
	fieldIndex := 2
	importer, err := NewImporter(&path, &skipHeader, &fieldIndex)
	if err != nil {
		t.Error(err)
	}
	data, err := importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}

	for i, v := range data {
		if v.Domain != sortedDomains[i] {
			t.Errorf("data not sorted properly. mismatch:\nhave: %v\nwant: %v", v.Domain, sortedDomains[i])
		}
	}
}

func TestImportInvalidPath(t *testing.T) {
	path := ""
	skipHeader := true
	fieldIndex := 2
	_, err := NewImporter(&path, &skipHeader, &fieldIndex)
	if err == nil {
		t.Error(err)
	}
}

func TestImportInvalidData(t *testing.T) {
	buf := bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	path := "./test_invalid_data.csv"
	skipHeader := true
	fieldIndex := 2
	importer, err := NewImporter(&path, &skipHeader, &fieldIndex)
	if err != nil {
		t.Error(err)
	}

	_, err = importer.ImportDomainData()

	out := buf.String()

	if !strings.Contains(out, "level=WARN msg=\"invalid email address\" email=\"\"") ||
		!strings.Contains(out, "level=WARN msg=\"invalid email address\" email=\"\"") {
		t.Error("invalid logging of invalid data", err)
	}

	if err != nil {
		t.Error("invalid data not caught", err)
	}
}

func BenchmarkImportDomainData(b *testing.B) {
	b.StopTimer()
	path := "./benchmark10k.csv"
	skipHeader := true
	fieldIndex := 2
	importer, err := NewImporter(&path, &skipHeader, &fieldIndex)
	if err != nil {
		b.Error(err)
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := importer.ImportDomainData(); err != nil {
			b.Error(err)
		}
	}
}
