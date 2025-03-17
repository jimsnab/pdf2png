package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

func pdfToPNG(pdfPath, outputDir string, dpi int) error {
	// Open the PDF document
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("error opening PDF: %v", err)
	}
	defer doc.Close()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Convert each page
	for i := 0; i < doc.NumPage(); i++ {
		// Extract image from page with specified DPI
		img, err := doc.ImageDPI(i, float64(dpi))
		if err != nil {
			return fmt.Errorf("error extracting image from page %d: %v", i+1, err)
		}

		// Create output file
		outputPath := filepath.Join(outputDir, fmt.Sprintf("page_%03d.png", i+1))
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating output file for page %d: %v", i+1, err)
		}

		// Encode image to PNG
		err = png.Encode(outFile, img)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("error encoding PNG for page %d: %v", i+1, err)
		}

		fmt.Printf("Converted page %d to %s\n", i+1, outputPath)
	}

	return nil
}

func printUsage() {
	fmt.Println("PDF to PNG Converter")
	fmt.Println("Usage: pdf2png <input-file> <output-dir> [-dpi=<resolution>]")
	fmt.Println("Arguments:")
	fmt.Println("  <input-file>    Input PDF file path (required)")
	fmt.Println("  <output-dir>    Output directory for PNG files (required)")
	fmt.Println("Flags:")
	fmt.Println("  -dpi            Resolution in DPI (default: 300)")
	fmt.Println("  -h, -help       Show this help message")
	fmt.Println("Note: Requires MuPDF installed on the system")
}

func main() {
	// Define optional flag
	dpi := flag.Int("dpi", 300, "Resolution in DPI (dots per inch)")
	help := flag.Bool("h", false, "Show help message")
	flag.BoolVar(help, "help", false, "Show help message") // Support both -h and -help

	// Parse flags
	flag.Parse()

	// Get positional arguments
	args := flag.Args()

	// Show help if requested or if arguments are invalid
	if *help || len(args) != 2 {
		printUsage()
		if !*help && len(args) != 2 {
			if len(args) < 2 {
				fmt.Println("\nError: missing required arguments")
			} else {
				fmt.Println("\nError: too many arguments")
			}
		}
		os.Exit(1)
	}

	inputFile := args[0]
	outputDir := args[1]

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: input file '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	// Perform conversion
	fmt.Printf("Converting %s to PNG files in %s directory (DPI: %d)...\n",
		inputFile, outputDir, *dpi)

	err := pdfToPNG(inputFile, outputDir, *dpi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Conversion completed successfully")
}