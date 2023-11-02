package nnvisualizer

import (
	"github.com/veandco/go-sdl2/sdl"
)

// InterpRatio:
func InterpRatio(arr []float64) []int {
	// constants
	minX := -10.0
	maxX := 10.0

	// interp between arr and the output (ratio)
	ratio := make([]int, len(arr))
	for k := 0; k < len(arr); k++ {
		value := arr[k]
		if value < minX {
			ratio[k] = 0
		} else if value > maxX {
			ratio[k] = 255
		} else {
			m := (255.0 - 0.0) / (maxX - minX)
			b := 0.0 - m*minX
			ratio[k] = int(m*float64(value) + b)
		}
	}

	return ratio
}

// ColorBlend: blend color by a weight
func ColorBlend(w int) sdl.Color {
	// blend between two colors
	endR, endG, endB := 0, 0, 255
	startR, startG, startB := 20, 20, 20

	// compute blended RGB values
	R := (1.0-float64(w)/255.0)*float64(startR) + (float64(w)/255.0)*float64(endR)
	G := (1.0-float64(w)/255.0)*float64(startG) + (float64(w)/255.0)*float64(endG)
	B := (1.0-float64(w)/255.0)*float64(startB) + (float64(w)/255.0)*float64(endB)

	return sdl.Color{
		R: uint8(R),
		G: uint8(G),
		B: uint8(B),
		A: uint8(255),
	}
}
