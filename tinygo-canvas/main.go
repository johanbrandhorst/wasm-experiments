// +build js,wasm

package main

import (
	"math"
	"math/rand"
	"strconv"
	"syscall/js"
)

var (
	width    float64
	height   float64
	mousePos [2]float64
	ctx, doc js.Value
	canvasEl js.Value
	dt       DotThing

	tmark     float64
	markCount = 0
	tdiffSum  float64
)

func main() {
	// Init Canvas stuff
	doc = js.Global().Get("document")
	canvasEl = doc.Call("getElementById", "mycanvas")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Call("setAttribute", "width", width)
	canvasEl.Call("setAttribute", "height", height)
	ctx = canvasEl.Call("getContext", "2d")

	// Set up dot thing
	dt = DotThing{speed: 160}
	dt.SetNDots(100)
	dt.lines = false

	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}

// * Handlers for JS callback functions *

//go:export speedInput
func speedInput(fval float64) {
	dt.speed = fval
}

//go:export countChange
func countChange(intVal int) {
	dt.SetNDots(intVal)
}

//go:export moveHandler
func moveHandler(cx int, cy int) {
	mousePos[0] = float64(cx)
	mousePos[1] = float64(cy)
}

//go:export renderFrame
func renderFrame(now float64) {
	tdiff := now - tmark
	tdiffSum += now - tmark
	markCount++
	if markCount > 10 {
		doc.Call("getElementById", "fps").Set("innerHTML", "FPS: "+strconv.FormatFloat(1000/(tdiffSum/float64(markCount)), 'f', 1, 64))
		tdiffSum, markCount = 0, 0
	}
	tmark = now

	// Pool window size to handle resize
	curBodyW := doc.Get("body").Get("clientWidth").Float()
	curBodyH := doc.Get("body").Get("clientHeight").Float()
	if curBodyW != width || curBodyH != height {
		width, height = curBodyW, curBodyH
		canvasEl.Set("width", width)
		canvasEl.Set("height", height)
	}
	dt.Update(tdiff / 1000)
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}

// DotThing manager
type DotThing struct {
	dots  []*Dot
	lines bool
	speed float64
}

// Update updates the dot positions and draws
func (dt *DotThing) Update(dtTime float64) {
	if dt.dots == nil {
		return
	}
	ctx.Call("clearRect", 0, 0, width, height)

	// Update
	for i, dot := range dt.dots {
		if dot.pos[0] < dot.size {
			dot.pos[0] = dot.size
			dot.dir[0] *= -1
		}
		if dot.pos[0] > width-dot.size {
			dot.pos[0] = width - dot.size
			dot.dir[0] *= -1
		}

		if dot.pos[1] < dot.size {
			dot.pos[1] = dot.size
			dot.dir[1] *= -1
		}

		if dot.pos[1] > height-dot.size {
			dot.pos[1] = height - dot.size
			dot.dir[1] *= -1
		}

		mdx := mousePos[0] - dot.pos[0]
		mdy := mousePos[1] - dot.pos[1]
		d := math.Sqrt(mdx*mdx + mdy*mdy)
		if d < 200 {
			dInv := 1 - d/200
			dot.dir[0] += (-mdx / d) * dInv * 8
			dot.dir[1] += (-mdy / d) * dInv * 8
		}
		for j, dot2 := range dt.dots {
			if i == j {
				continue
			}
			mx := dot2.pos[0] - dot.pos[0]
			my := dot2.pos[1] - dot.pos[1]
			d := math.Sqrt(mx*mx + my*my)
			if d < 100 {
				dInv := 1 - d/100
				dot.dir[0] += (-mx / d) * dInv
				dot.dir[1] += (-my / d) * dInv
			}
		}
		dot.dir[0] *= 0.1 //friction
		dot.dir[1] *= 0.1 //friction

		dot.pos[0] += dot.dir[0] * dt.speed * dtTime * 10
		dot.pos[1] += dot.dir[1] * dt.speed * dtTime * 10

		ctx.Set("globalAlpha", 0.5)
		ctx.Call("beginPath")
		hexCol := hexFormat(dot.color)
		ctx.Set("fillStyle", hexCol)
		ctx.Set("strokeStyle", hexCol)
		ctx.Set("lineWidth", dot.size)
		ctx.Call("arc", dot.pos[0], dot.pos[1], dot.size, 0, 2*math.Pi)
		ctx.Call("fill")
	}
}

// SetNDots reinitializes dots with n size
func (dt *DotThing) SetNDots(n int) {
	dt.dots = make([]*Dot, n)
	for i := 0; i < n; i++ {
		dt.dots[i] = &Dot{
			pos: [2]float64{
				rand.Float64() * width,
				rand.Float64() * height,
			},
			dir: [2]float64{
				rand.NormFloat64(),
				rand.NormFloat64(),
			},
			color: uint32(rand.Intn(0xFFFFFF)),
			size:  10,
		}
	}
}

// Adapted from the Go source: https://github.com/golang/go/blob/4ce6a8e89668b87dce67e2f55802903d6eb9110a/src/fmt/format.go#L248-L252
func hexFormat(u uint32) string {
	digits := "0123456789abcdefx"
	buf := make([]uint8, 6)
	i := len(buf)
	for u >= 16 {
		i--
		buf[i] = digits[u&0xF]
		u >>= 4
	}
	i--
	buf[i] = digits[u]
	return "#" + string(buf)
}

// Dot represents a dot ...
type Dot struct {
	pos   [2]float64
	dir   [2]float64
	color uint32
	size  float64
}
