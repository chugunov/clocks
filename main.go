package main

import (
	"flag"
	"log"
)

func initFlags() (string, string) {
	var inputFile string
	var outputFile string

	flag.StringVar(
		&inputFile, "i", "",
		`Path to the input file containing the processes sequence. Example format:
p0: s1 r1 l r1
p1: s0 s2 r0 l s2 s0 l r2
p2: l s1 r1 r1`,
	)
	flag.StringVar(
		&outputFile, "o", "space_time_diagram.svg",
		`Path to the output file.
Supported extensions are: .eps, .jpg|.jpeg, .pdf, .png, .svg, .tex, .tif|.tiff`,
	)
	flag.Parse()

	if inputFile == "" {
		log.Fatal("Input file path must be provided with -i flag")
	}

	return inputFile, outputFile
}

func main() {
	inputFile, outputFile := initFlags()
	sim := NewSimulator(inputFile)
	err := sim.Simulate()
	if err != nil {
		log.Fatalf("Error during simulation: %e", err)
	}
	err = sim.Draw(outputFile)
	if err != nil {
		log.Fatalf("Error during render space-time diagram: %e", err)
	}
	log.Printf("Space-time diagram saved as %s\n", outputFile)
}
