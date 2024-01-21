package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	//Open the file
	file, err := os.Open(filename)
	//Check the potentiel error
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//Create PBM variable
	pbm := &PBM{}
	//Lines of the image
	line := 0
	//Loop each line
	for scanner.Scan() {
		text := scanner.Text()
		//Ignore comments ect....
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		if pbm.magicNumber == "" {
			//Make sure pbm.magicNumber contains text
			pbm.magicNumber = strings.TrimSpace(text)
		} else if pbm.width == 0 {
			//Get the width and height of the pbm
			fmt.Sscanf(text, "%d %d", &pbm.width, &pbm.height)
			//Initialize the pbm.data matrix variable by creating the correct amount and size of arrays in an array
			pbm.data = make([][]bool, pbm.height)
			for i := range pbm.data {
				pbm.data[i] = make([]bool, pbm.width)
			}
		} else {
			if pbm.magicNumber == "P1" {
				//Each word in the text string is printed on a new line
				test := strings.Fields(text)
				//Loop through the string[]
				for i := 0; i < pbm.width; i++ {
					//If the given string == "1", then it is stored as true or else as false
					pbm.data[line][i] = (test[i] == "1")
				}
				line++
			} else if pbm.magicNumber == "P4" {
				//Calculate the expected number of bytes per row
				expectedBytesPerRow := (pbm.width + 7) / 8
				totalExpectedBytes := expectedBytesPerRow * pbm.height
				allPixelData := make([]byte, totalExpectedBytes)
				//Reads the file content
				fileContent, err := os.ReadFile(filename)
				if err != nil {
					return nil, fmt.Errorf("couldn't read file: %v", err)
				}
				//Extracts the necessary pixel data
				copy(allPixelData, fileContent[len(fileContent)-totalExpectedBytes:])
				//Process the data to fill the pixel array of pbm.data
				byteIndex := 0
				for y := 0; y < pbm.height; y++ {
					for x := 0; x < pbm.width; x++ {
						if x%8 == 0 && x != 0 {
							byteIndex++
						}
						pbm.data[y][x] = (allPixelData[byteIndex]>>(7-(x%8)))&1 != 0
					}
					byteIndex++
				}
				break
			}
		}
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	//Return size
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	//Return value a pixel
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	//Define a new value pixel
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", pbm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	_, err = fmt.Fprintf(file, "%d %d\n", &pbm.width, &pbm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	for _, row := range pbm.data {
		for _, pixel := range row {
			if pbm.magicNumber == "P1" {
				if pixel {
					_, err = fmt.Fprint(file, "1 ")
				} else {
					_, err = fmt.Fprint(file, "0 ")
				}
			} else if pbm.magicNumber == "P4" {
				if pixel {
					_, err = fmt.Fprint(file, "1")
				} else {
					_, err = fmt.Fprint(file, "0")
				}
			}
			if err != nil {
				return fmt.Errorf("error writing data: %v", err)
			}
		}
		if pbm.magicNumber == "P1" {
			_, err = fmt.Fprintln(file)
			if err != nil {
				return fmt.Errorf("error writing data: %v", err)
			}
		}
	}

	fmt.Printf("Image sauvegardée avec succès dans %s\n", filename)

	return nil
}
