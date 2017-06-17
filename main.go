// Open an OpenGl window and display a rectangle using a OpenGl GraphicContext
package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
)

const (
	numCircles = 20	//number of moving circles
)

var (
	random *rand.Rand
	width, height = 512, 512
	floatwidth, floatheight = float64(width), float64(height)
)

type circle struct {
	xpos, ypos, radius, xdelta, ydelta float64
	intensity int
}

var circles [numCircles]circle

// called when the window is reshaped/resized
func reshape(window *glfw.Window, w, h int) {
	gl.ClearColor(0, 0, 0, 1)
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin up to upper-left corner. */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
	floatwidth, floatheight = float64(width), float64(height)
	initializeCircles()
}

func init() {
	// locks to a particular thread
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(width, height, "Radar", nil, nil)
	if err != nil {
		panic(err)
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	window.MakeContextCurrent()
	// set callback for when a window is resized
	window.SetSizeCallback(reshape)
	// sets callback for when a key is pressed (not keyboard layout dependent), repeated or released
	window.SetKeyCallback(onKey)
	// sets callback for when a key is input (keyboard layout dependent)
	window.SetCharCallback(onChar)

	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	reshape(window, width, height)

	// main loop
	for !window.ShouldClose() {
		drawContents(window)
		//window.SwapBuffers()
		glfw.PollEvents()
		//time.Sleep(2 * time.Second)
	}
}

// called when a character is input
func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

// called when a key is pressed, repeated or released
func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// quits on escape or letter "q" being pressed; uses physical keyboard layout (ignores dvorak on linux)
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

// initializes position, movement delta, and brightness for moving circles
func initializeCircles() {
	for i := 0; i < numCircles; i++ {
		c := circle{
			xpos: (floatwidth / 2) + (random.Float64() * float64(direction()) * floatwidth * 0.5),
			ypos: (floatheight / 2) + (random.Float64() * float64(direction()) * floatheight * 0.5),
			radius: random.Float64() * floatwidth * 0.03,
			xdelta: random.Float64() * float64(direction()) * floatwidth * 0.01,
			ydelta: random.Float64() * float64(direction()) * floatheight * 0.01,
			intensity: 255,
		}
		circles[i] = c
	}
}

// returns a random +1 or -1 result; can be used for a random left/right or up/down direction
func direction() int {
	if random.Intn(2) % 2 == 0 {
		return 1 
	} else {
		return -1
	}
}

// updates moving circles
func updateCircles(angle int) {
	if angle % 359 == 0 {
		for i := range circles {
			circle := circles[i]
			circle.xpos = circle.xpos + circle.xdelta
			circle.ypos = circle.ypos + circle.ydelta
			circle.intensity = 255
			circles[i] = circle
		}
	}
	if angle % 4 == 0 {
		for i := range circles {
			circle := circles[i]
			circle.intensity = circle.intensity - 2
			circles[i] = circle
		}
	}
}

// draws the contents of the window
func drawContents(w *glfw.Window) {
	for i := 0; i < 360; i++ {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LineWidth(1)

		gc := draw2dgl.NewGraphicContext(width, height)

		drawRadials(gc, floatwidth, floatheight, floatwidth, floatheight)
		drawConcentricCircles(gc, floatwidth, floatheight, floatwidth, floatheight)
		drawSweep(gc, floatwidth, floatheight, floatwidth, floatheight, i)

		for j := range circles {
			drawMovingCircle(gc, circles[j].xpos, circles[j].ypos, circles[j].radius, circles[j].radius, circles[j].intensity)
		}
		updateCircles(i)

		drawMaskingCircle(gc, floatwidth, floatheight, floatwidth, floatheight)

		gl.Flush() /* single buffered, so needs a flush. */
		time.Sleep(2 * time.Millisecond)
		w.SwapBuffers()
		//time.Sleep(2 * time.Millisecond)
	}
}

// draws radial lines in radar
func drawRadials(gc *draw2dgl.GraphicContext, x, y, width, height float64) {
	gc.Save()
	gc.Translate(x/2, y/2)
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})
	for i := 0.0; i < 360; i = i + 45 { // go from 0 to 360 degrees in 45 degree steps
		gc.Save()                        // keep rotations temporary
		gc.Rotate(i * (math.Pi / 180.0)) // rotate by degrees on stack from 'for'
		gc.MoveTo(0, 0)
		if width > height {
			gc.LineTo(height * 0.5, 0)
		} else {
			gc.LineTo(width * 0.5, 0)
		}
		gc.Stroke()
		gc.Restore()
	}
	gc.Restore()
}

// draws concentric circles in radar
func drawConcentricCircles(gc *draw2dgl.GraphicContext, xc, yc, width, height float64) {
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})

	var radius float64
	if width > height {
		radius = height/2
	} else {
		radius = width/2
	}
	xc = width/2
	yc = height/2
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0)     /* clockwise in radians           */
	gc.SetLineCap(draw2d.ButtCap)
	for i := 1.0; i > 0; i = i - 0.3 { // reduction factor for concentric circles
		gc.MoveTo(xc + math.Cos(startAngle) * radius, yc + math.Sin(startAngle) * radius)
		gc.ArcTo(xc, yc, radius * i, radius * i, startAngle, sweepAngle)
		gc.Stroke()
	}
}

// draws radar sweep
func drawSweep(gc *draw2dgl.GraphicContext, x, y, width, height float64, angle int) {
	var length float64
	if width > height {
		length = height/2
	} else {
		length = width/2
	}
	for i := 0; i < 60; i++ {
		gc.Save()
		gc.Translate(x/2, y/2)
		gc.SetLineWidth(1)
		gc.SetStrokeColor(color.RGBA{0, uint8(250 - 3 * i), 0, 0xff})
		gc.Rotate((float64(angle) - float64(i) * 0.5) * (math.Pi / 180.0))
		gc.MoveTo(0, 0)
		gc.LineTo(length, 0)
		gc.Stroke()
		gc.Restore()
	}
}

// draws circles that move
func drawMovingCircle(gc *draw2dgl.GraphicContext, xc, yc, width, height float64, intensity int) {
	gc.SetLineWidth(4)
	gc.SetStrokeColor(color.RGBA{0, uint8(intensity), 0, 0xff})
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0)     /* clockwise in radians           */
	gc.SetLineCap(draw2d.ButtCap)
	gc.MoveTo(xc + math.Cos(startAngle) * width, yc + math.Sin(startAngle) * height)
	gc.ArcTo(xc, yc, width, height, startAngle, sweepAngle)
	gc.Stroke()
}

// draws black circle to mask moving circles that move outside of the outer radar circle
func drawMaskingCircle(gc *draw2dgl.GraphicContext, xc, yc, width, height float64) {
	var radius float64
	if width > height {
		radius = height/2
	} else {
		radius = width/2
	}
	xc = width/2
	yc = height/2
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0)     /* clockwise in radians           */
	gc.SetLineWidth(radius)
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
	gc.SetLineCap(draw2d.ButtCap)
	gc.MoveTo(xc + math.Cos(startAngle) * radius, yc + math.Sin(startAngle) * radius)
	gc.ArcTo(xc, yc, radius + (radius/2) + 1, radius + (radius/2) + 1, startAngle, sweepAngle)
	gc.Stroke()
}

