package main

import (
	"bufio"
	"fmt"
	"os"
)

type Image struct {
	EnhancementAlgorithm []rune
	CurrentImage         [][]rune
	enhancementRound     uint
}

func (image *Image) Pad(filled bool) {
	paddedImage := make([][]rune, len(image.CurrentImage)+2)
	for i := range paddedImage {
		// We assume that all lines have the same width in the image.
		paddedImage[i] = make([]rune, len(image.CurrentImage[0])+2)
		if i >= 1 && i < len(paddedImage)-1 {
			for j := range paddedImage[i] {
				if j >= 1 && j < len(paddedImage[i])-1 {
					paddedImage[i][j] = image.CurrentImage[i-1][j-1]
				} else {
					if filled {
						paddedImage[i][j] = '#'
					} else {
						paddedImage[i][j] = '.'
					}
				}
			}
		} else {
			for j := range paddedImage[i] {
				if filled {
					paddedImage[i][j] = '#'
				} else {
					paddedImage[i][j] = '.'
				}
			}
		}
	}

	image.CurrentImage = paddedImage
}

func (image *Image) Enhance() {
	image.enhancementRound++
	// We only check if the 0 bit for image enhancement is filled or not.
	// Technically this shouldn't be enough, but for this task it is, otherwise
	// it would be impossible to solve.
	filedPadding := image.enhancementRound%2 == 0 && image.EnhancementAlgorithm[0] == '#'
	// we adjust the padding from the previous round as well
	for i := 0; i < len(image.CurrentImage); i++ {
		for j := 0; j < len(image.CurrentImage[i]); j++ {
			if i == 0 || j == 0 || i == len(image.CurrentImage)-1 || j == len(image.CurrentImage[i])-1 {
				if filedPadding {
					image.CurrentImage[i][j] = '#'
				} else {
					image.CurrentImage[i][j] = '.'
				}
			}
		}
	}
	if filedPadding {
		image.Pad(true)
	} else {
		image.Pad(false)
	}

	enhancedImage := make([][]rune, len(image.CurrentImage))
	for i := range enhancedImage {
		enhancedImage[i] = make([]rune, len(image.CurrentImage[i]))
		for j := range enhancedImage[i] {
			enhancedImage[i][j] = '.'
		}
	}
	for i := 1; i < len(enhancedImage)-1; i++ {
		for j := 1; j < len(enhancedImage[i])-1; j++ {
			enhancementCode := make([]rune, 0, 9)
			enhancementCode = append(enhancementCode, image.CurrentImage[i-1][j-1:j+2]...)
			enhancementCode = append(enhancementCode, image.CurrentImage[i][j-1:j+2]...)
			enhancementCode = append(enhancementCode, image.CurrentImage[i+1][j-1:j+2]...)

			index := uint(0)
			for _, r := range enhancementCode {
				index = index << 1
				if r == '#' {
					index |= 1
				}
			}
			enhancedImage[i][j] = image.EnhancementAlgorithm[index]
		}
	}
	image.CurrentImage = enhancedImage
}

func StringToRunes(input string) []rune {
	data := make([]rune, 0, len(input))
	for _, r := range input {
		data = append(data, r)
	}
	return data
}

func ParseInput(lines []string) Image {
	image := Image{}
	image.EnhancementAlgorithm = StringToRunes(lines[0])

	for _, line := range lines[1:] {
		if len(line) == 0 {
			continue
		}
		image.CurrentImage = append(image.CurrentImage, StringToRunes(line))
	}
	image.Pad(false)
	return image
}

func EnhanceAndCount(image Image, times uint) uint {
	for count := uint(0); count < times; count++ {
		image.Enhance()
	}
	lit := uint(0)
	for _, line := range image.CurrentImage {
		for _, r := range line {
			if r == '#' {
				lit++
			}
		}
	}
	return lit
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	image_one := ParseInput(lines)
	image_two := ParseInput(lines) // easier than to write a copy
	
	fmt.Println()
	fmt.Printf("The number of lit pixels after 2 enhancement rounds is: %v\n", EnhanceAndCount(image_one, 2))
	fmt.Printf("The number of lit pixels after 50 enhancement rounds is: %v\n", EnhanceAndCount(image_two, 50))
}
