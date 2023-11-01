package nnvisualizer

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// WinProp: properties of display window
type WinProp struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Graph struct {
	Props    WinProp
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

type NNLayer struct {
	Neurons     int
	Activations []int
}

type NNGraph struct {
	NumberOfLayers int
	NNLayers       []NNLayer
	WinGraph       *Graph
}

// Init: start sdl library
func Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}
	return nil
}

// NewNNGraph:
func NewNNGraph(n int, g *Graph) (*NNGraph, error) {
	// number of layers greater than 0
	if n <= 0 {
		return nil, fmt.Errorf("number of layers must be > 0. got %d", n)
	}

	return &NNGraph{
		NumberOfLayers: n,
		NNLayers:       make([]NNLayer, n),
		WinGraph:       g,
	}, nil
}

func RenderCircle(r *sdl.Renderer, x, y, radius int, c sdl.Color) {
	for i := 0; i < 2*radius; i++ {
		for j := 0; j < 2*radius; j++ {
			dx := radius - i
			dy := radius - j
			if (dx*dx + dy*dy) <= radius*radius {
				_ = r.SetDrawColor(c.R, c.G, c.B, c.A)
				_ = r.DrawPoint(int32(x+dx), int32(y+dy))
			}
		}
	}
}

// NewGraph: creates sdl window and renderer
func NewGraph(wp WinProp) (*Graph, error) {
	w, err := sdl.CreateWindow(
		"neural network",
		int32(wp.X),
		int32(wp.Y),
		int32(wp.Width),
		int32(wp.Height),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, err
	}

	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	if err := r.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return nil, err
	}

	return &Graph{
		Props:    wp,
		Window:   w,
		Renderer: r,
	}, nil
}

func (n *NNGraph) RenderNetwork() error {
	// shorthand for sdl2 renderer
	r := n.WinGraph.Renderer

	// background
	if err := r.SetDrawColor(20, 40, 50, 255); err != nil {
		return err
	}

	if err := r.Clear(); err != nil {
		return err
	}

	// layer width & offset in X direction
	layerWidth := n.WinGraph.Props.Width / len(n.NNLayers)
	xOffset := layerWidth / 2

	// loop over layers
	for i := 0; i < len(n.NNLayers); i++ {
		layer := n.NNLayers[i]
		for j := 0; j < layer.Neurons; j++ {
			// positions of neurons (x,y)
			x := xOffset + i*layerWidth
			neuronHeight := 1 + n.WinGraph.Props.Height/layer.Neurons
			yOffset := neuronHeight / 2
			y := yOffset + j*neuronHeight

			// connect neurons to next layer
			if i < len(n.NNLayers)-1 {
				nextLayer := n.NNLayers[i+1]
				for k := 0; k < nextLayer.Neurons; k++ {
					xNext := xOffset + (i+1)*layerWidth
					neuronHeightNext := 1 + n.WinGraph.Props.Height/nextLayer.Neurons
					yOffsetNext := neuronHeightNext / 2
					yNext := yOffsetNext + k*neuronHeightNext

					_ = r.SetDrawColor(0x55, 0x55, 0x55, 80) // color of links
					_ = r.DrawLine(int32(x), int32(y), int32(xNext), int32(yNext))
				}
			}

			// draw neurons
			bStrength := layer.Activations[j]
			c := genColorGradient(bStrength)
			neuronSize := 1 + (100 / layer.Neurons)
			RenderCircle(r, x, y, neuronSize, c)
		}
	}

	r.Present()
	return nil
}

func genColorGradient(num int) sdl.Color {
	startR, startG := 15, 10
	w := float64(num) / 255

	r := float64(startR) * (1 - w)
	g := float64(startG) * (1 - w)
	b := float64(255) * w

	return sdl.Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(255),
	}
}

// Close: call upon closing the display graph
func (g *Graph) Close() {
	g.Window.Destroy()
	g.Renderer.Destroy()
	sdl.Quit()
}
