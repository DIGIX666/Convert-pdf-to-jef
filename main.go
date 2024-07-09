package main

import (
	"fmt"
)

func main() {
	// read SVG
	svg, err := readSVG("design.svg")
	if err != nil {
		fmt.Println("Error reading SVG:", err)
		return
	}

	// Convert SVG to JEF points (embroidery points)
	points := convertSVGToJEF(svg)

	// Write JEF
	err = writeJEF("design.jef", points)
	if err != nil {
		fmt.Println("Error writing JEF:", err)
		return
	}

	fmt.Println("Conversion completed successfully")
}
