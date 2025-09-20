package main

import (
	"flag"
	"fmt"
	"log"
	"parser/exporter"
	"parser/importer"
	"strings"
)

type Options struct {
	path       *string
	skipHeader *bool
	fieldIndex *int
	outFile    *string
}

func readOptions() (*Options, error) {
	opts := &Options{}
	opts.path = flag.String("path", "./customers.csv", "Path to the file with customer data")
	opts.skipHeader = flag.Bool("skipHeader", true, "Skip header of csv file")
	opts.fieldIndex = flag.Int("field", 0, "Index of email field in csv file")
	opts.outFile = flag.String("out", "./output.csv", "Optional: output file path. If empty program will output results to the terminal")
	flag.Parse()

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return opts, nil
}

func (o *Options) Validate() error {
	if o.path == nil || *o.path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if o.fieldIndex == nil || *o.fieldIndex < 0 {
		return fmt.Errorf("field index must be >= 0, got %d", *o.fieldIndex)
	}
	if o.outFile != nil && *o.outFile != "" && !strings.HasSuffix(*o.outFile, ".csv") {
		return fmt.Errorf("output file must have .csv extension")
	}
	return nil
}

func main() {
	opts, err := readOptions()
	if err != nil {
		log.Fatal(err)
	}

	dataImporter, err := importer.NewImporter(opts.path, opts.skipHeader, opts.fieldIndex)
	if err != nil {
		log.Fatal(err)
	}

	domainData, err := dataImporter.ImportDomainData()
	if err != nil {
		log.Fatal(err)
	}

	dataExporter, err := exporter.NewExporter(opts.outFile)
	if err := dataExporter.ExportData(domainData); err != nil {
		log.Fatal(err)
	}
}
