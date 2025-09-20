package exporter

import (
	"fmt"
	"parser/data"
	"parser/importer"
	"testing"
)

func TestExportData(t *testing.T) {
	path := "./test_output.csv"
	domainData := []data.CustomerDomainData{
		{
			Domain:           "livejournal.com",
			CustomerQuantity: 12,
		},
		{
			Domain:           "microsoft.com",
			CustomerQuantity: 22,
		},
		{
			Domain:           "newsvine.com",
			CustomerQuantity: 15,
		},
		{
			Domain:           "pinteres.uk",
			CustomerQuantity: 10,
		},
		{
			Domain:           "yandex.ru",
			CustomerQuantity: 43,
		},
	}
	exporter, err := NewExporter(&path)
	if err != nil {
		t.Fatal(err)
	}

	err = exporter.ExportData(domainData)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportEmptyData(t *testing.T) {
	path := "./test_output.csv"
	exporter, err := NewExporter(&path)
	if err != nil {
		t.Fatal(err)
	}

	err = exporter.ExportData(nil)
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func BenchmarkImportDomainData(b *testing.B) {
	b.StopTimer()
	dir := b.TempDir()
	path := fmt.Sprintf("%s/test_output.csv", dir)
	dataPath := "../customerimporter/benchmark10k.csv"
	skipHeader := true
	fieldIndex := 2
	importer, err := importer.NewImporter(&dataPath, &skipHeader, &fieldIndex)
	if err != nil {
		b.Fatal(err)
	}
	data, err := importer.ImportDomainData()
	if err != nil {
		b.Error(err)
	}
	exporter, err := NewExporter(&path)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := exporter.ExportData(data); err != nil {
			b.Fatal(err)
		}
	}
}
