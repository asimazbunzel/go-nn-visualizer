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
	Synapse     [][]int
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

// RenderCircle: render a full circle
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

// RenderNetwork:
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

	// layer width, offset in X direction & initial x position
	layerWidth := n.WinGraph.Props.Width / len(n.NNLayers)
	xOffset := layerWidth / 2
	xBegin, yBegin := 1, 1

	// loop over layers
	for i := 0; i < len(n.NNLayers); i++ {
		// shortcut
		layer := n.NNLayers[i]

		// loop over neurons at a fixed layer
		for j := 0; j < layer.Neurons; j++ {
			// positions of neurons (x,y)
			x := xBegin + xOffset + i*layerWidth
			neuronHeight := n.WinGraph.Props.Height / layer.Neurons
			yOffset := neuronHeight / 2
			y := yBegin + yOffset + j*neuronHeight

			// connect neurons to next layer (synapse)
			if i < len(n.NNLayers)-1 { // no need to do this for the output layer
				// shortcut
				nextLayer := n.NNLayers[i+1]

				// link each neuron in the layer with neurons in the next layer
				for k := 0; k < nextLayer.Neurons; k++ {
					// positions of neurons in next layer (xNext, yNext)
					xNext := xBegin + xOffset + (i+1)*layerWidth
					neuronHeightNext := n.WinGraph.Props.Height / nextLayer.Neurons
					yOffsetNext := neuronHeightNext / 2
					yNext := yBegin + yOffsetNext + k*neuronHeightNext

					_ = r.SetDrawColor(0x55, 0x55, 0x55, 80) // color of links (TODO: change according to synapse values)
					_ = r.DrawLine(int32(x), int32(y), int32(xNext), int32(yNext))
				}
			}

			// draw neurons
			bStrength := layer.Activations[j]
			c := ColorBlend(bStrength)
			neuronSize := 1 + (100 / layer.Neurons)
			RenderCircle(r, x, y, neuronSize, c)
		}
	}

	r.Present()
	return nil
}

// Close: call upon closing the display graph
func (g *Graph) Close() {
	g.Window.Destroy()
	g.Renderer.Destroy()
	sdl.Quit()
}
