package exporter

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"parser/data"
	"strconv"
)

type Exporter struct {
	outputPath *string
	headers    []string
}

// NewExporter returns a new Exporter that writes customer domain data to specified file.
func NewExporter(outputPath *string) (*Exporter, error) {

	return &Exporter{
		outputPath: outputPath,
		headers:    []string{"domain", "number_of_customers"},
	}, nil
}

// ExportData if outputPath is empty prints data to terminal otherwise to CSV file.
func (ex Exporter) ExportData(data []data.CustomerDomainData) error {
	if data == nil {
		return fmt.Errorf("provided data is empty (nil)")
	}

	if *ex.outputPath == "" {
		ex.printData(data)
	} else {
		err := ex.exportCsv(data)
		if err != nil {
			return err
		}
	}

	return nil
}

// printData prints data to terminal
func (ex Exporter) printData(data []data.CustomerDomainData) {
	fmt.Println("domain,number_of_customers")
	for _, v := range data {
		fmt.Printf("%s,%v\n", v.Domain, v.CustomerQuantity)
	}
}

// exportCsv exports data to CSV file. If file already exists, it will be truncated.
func (ex Exporter) exportCsv(data []data.CustomerDomainData) error {
	outputFile, err := os.Create(*ex.outputPath)
	if err != nil {
		return err
	}
	slog.Info("opened file", "path", *ex.outputPath)
	defer func() {
		if err := outputFile.Close(); err != nil {
			slog.Error("failed to close file", "err", err)
		} else {
			slog.Info("closed file", "path", *ex.outputPath)
		}
	}()
	csvWriter := csv.NewWriter(outputFile)
	defer func() {
		csvWriter.Flush()
		if err := csvWriter.Error(); err != nil {
			slog.Error("failed to flush CSV writer", "err", err)
		} else {
			slog.Info("flushed CSV writer")
		}
	}()
	if err := csvWriter.Write(ex.headers); err != nil {
		return err
	}
	for _, v := range data {
		pair := []string{v.Domain, strconv.FormatUint(v.CustomerQuantity, 10)}
		if err := csvWriter.Write(pair); err != nil {
			return err
		}
	}

	return nil
}
