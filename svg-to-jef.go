package main

import (
	"encoding/binary"
	"encoding/xml"
	"os"
	"strconv"
	"strings"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Paths   []Path   `xml:"path"`
}

type Path struct {
	D string `xml:"d,attr"`
}

func readSVG(filename string) (SVG, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return SVG{}, err
	}

	var svg SVG
	err = xml.Unmarshal(data, &svg)
	if err != nil {
		return SVG{}, err
	}

	return svg, nil
}

// Point represents an embroidery point
type Point struct {
	X int16
	Y int16
}

// Convert SVG to JEF points (embroidery points)
func convertSVGToJEF(svg SVG) []Point {
	var points []Point

	for _, path := range svg.Paths {
		d := path.D
		commands := strings.Fields(d)

		var x, y int16
		for _, cmd := range commands {
			switch cmd[0] {
			case 'M', 'L':
				coords := strings.Split(cmd[1:], ",")
				x = parseCoord(coords[0])
				y = parseCoord(coords[1])
				points = append(points, Point{X: x, Y: y})
			case 'm', 'l':
				coords := strings.Split(cmd[1:], ",")
				x += parseCoord(coords[0])
				y += parseCoord(coords[1])
				points = append(points, Point{X: x, Y: y})
			default:
			}
		}
	}

	return points
}

// Convert SVG coordinates to embroidery stitches (adjust as required)
// Here we simplify things by converting directly to int16
func parseCoord(coord string) int16 {
	value, _ := strconv.ParseInt(coord, 10, 16)
	return int16(value)
}

func writeJEF(filename string, points []Point) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var header [128]byte
	binary.Write(file, binary.LittleEndian, header)

	for _, point := range points {
		binary.Write(file, binary.LittleEndian, point.X)
		binary.Write(file, binary.LittleEndian, point.Y)
	}

	return nil
}
