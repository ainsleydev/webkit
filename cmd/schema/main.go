package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/printer"
)

func main() {
	outputPath := flag.String("output", ".webkit/schema.json", "Path to output schema file")
	stdout := flag.Bool("stdout", false, "Output schema to stdout instead of file")
	flag.Parse()

	p := printer.New(os.Stdout)

	if err := run(context.Background(), p, *outputPath, *stdout); err != nil {
		p.Error(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}

func run(_ context.Context, p *printer.Console, outputPath string, stdout bool) error {
	p.Info("Generating JSON schema...")

	schemaData, err := appdef.GenerateSchema()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	if stdout {
		fmt.Println(string(schemaData))
		return nil
	}

	dir := filepath.Dir(outputPath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	if err = os.WriteFile(outputPath, schemaData, 0644); err != nil {
		return fmt.Errorf("writing schema file: %w", err)
	}

	p.Success("Schema generated successfully at: " + outputPath)

	return nil
}
