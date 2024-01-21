package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

func ReadPGM(filename string) (*PGM, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	//Check the potentiel error
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Create a buffered reader for efficient reading
	reader := bufio.NewReader(file)
	// Read magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}
	//Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width and height must be positive")
	}
	//Read max value
	maxValue, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading max value: %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	var max uint8
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid max value: %v", err)
	}
	// Read image data
	data := make([][]uint8, height)
	if magicNumber == "P2" {
		// Read P2 format
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]uint8, width)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				var pixelValue uint8
				_, err := fmt.Sscanf(field, "%d", &pixelValue)
				if err != nil {
					return nil, fmt.Errorf("error parsing pixel value at row %d, column %d: %v", y, x, err)
				}
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	}
	// Return the PGM struct
	return &PGM{data, width, height, magicNumber, max}, nil
}

// Return size
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// Return value a pixel
func (pgm *PGM) At(x, y int) uint8 {
	if x >= 0 && x < pgm.width && y >= 0 && y < pgm.height {
		return pgm.data[y][x]
	}
	return 0
}

// Define a new value pixel
func (pgm *PGM) Set(x, y int, value uint8) {
	if x >= 0 && x < pgm.width && y >= 0 && y < pgm.height {
		pgm.data[y][x] = value
	}
}

func (pgm *PGM) Invert() {
	maxUint8 := uint8(pgm.max)
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			//Invert pixel value
			pgm.data[y][x] = maxUint8 - pgm.data[y][x]
		}
	}
}

func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		left := 0
		right := pgm.width - 1
		//Invert pixel values from left to right
		for left < right {
			pgm.data[y][left], pgm.data[y][right] = pgm.data[y][right], pgm.data[y][left]
			left++
			right--
		}
	}
}
