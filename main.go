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
	numCircles = 20 //number of moving circles
)

var (
	random        *rand.Rand
	width, height = 512.0, 512.0
)

type Circle struct {
	xpos, ypos, radius, xdelta, ydelta float64
	brightness                         int //brightness of nircle drawn on screen
	pauseCounter                       int // how many ticks to wait before processing circle's position for potential movement
}

var circles [numCircles]Circle // collection of all moving circles

// Called when the window is reshaped/resized.
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
	width, height = float64(w), float64(h)
	initializeCircles()
}

func init() {
	// locks to a particular thread
	runtime.LockOSThread()
}

func main() {
	// for window resize callback
	windowWidth, windowHeight := int(width), int(height)

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(int(width), int(height), "Radar", nil, nil)
	if err != nil {
		panic(err)
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))

	window.MakeContextCurrent()

	// sets callback for when a window is resized
	window.SetSizeCallback(reshape)
	// sets callback for when a key is pressed (not keyboard layout dependent), repeated or released
	window.SetKeyCallback(onKey)
	// sets callback for when a key is input (keyboard layout dependent)
	window.SetCharCallback(onChar)

	// vertical synchronization
	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	// sets up initial window
	reshape(window, windowWidth, windowHeight)

	// main loop
	for !window.ShouldClose() {
		drawContents(window)
		glfw.PollEvents()
		//time.Sleep(2 * time.Second)
	}
}

// Called when a character is input.
func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

// Called when a key is pressed, repeated or released.
func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// quits on escape or letter "q" being pressed; uses physical keyboard layout (ignores dvorak on linux)
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}

// Initializes position, movement delta, and brightness for moving circles.
func initializeCircles() {
	for i := 0; i < numCircles; i++ {
		c := Circle{
			xpos:         (width / 2) + (random.Float64() * float64(direction()) * width * 0.5),
			ypos:         (height / 2) + (random.Float64() * float64(direction()) * height * 0.5),
			radius:       random.Float64() * width * 0.03,
			xdelta:       random.Float64() * float64(direction()) * width * 0.01,
			ydelta:       random.Float64() * float64(direction()) * height * 0.01,
			brightness:   255,
			pauseCounter: 0,
		}
		circles[i] = c
	}
}

// Returns a random +1 or -1 result. Can be used for a random left/right or up/down direction.
func direction() int {
	if random.Intn(2)%2 == 0 {
		return 1
	} else {
		return -1
	}
}

// Updates position and brightness of moving circles.
func updateCircles(angle int) {
	for i := range circles {
		circle := circles[i]
		if shouldMoveCircle(circle, angle) {
			circle.xpos = circle.xpos + circle.xdelta
			circle.ypos = circle.ypos + circle.ydelta
			circle.brightness = 255
			circle.pauseCounter = 50
		} else {
			circle.pauseCounter = circle.pauseCounter - 1
		}
		if angle%4 == 0 {
			circle.brightness = circle.brightness - 2
		}
		circles[i] = circle
	}
}

// Returns whether the given circle lines up with the sweep angle, and should be updated or not.
func shouldMoveCircle(c Circle, sweepAngle int) bool {
	// position that circle would move to, relative to the center of the radar
	xpos := (c.xpos + c.xdelta) - width/2
	ypos := (c.ypos + c.ydelta) - height/2
	// use arctan to find angle in radians; convert to degrees
	angle := math.Atan2(ypos, xpos) * 180 / math.Pi
	if angle < 0 {
		angle = 360 + angle
	}
	// the pause counter is used to prevent multiple consequtive moves for circles that only move a
	// small amount each time
	return int(angle) == sweepAngle && c.pauseCounter < 1
}

// Draws the contents of the window.
func drawContents(w *glfw.Window) {
	for i := 0; i < 360; i++ {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LineWidth(1)

		gc := draw2dgl.NewGraphicContext(int(width), int(height))

		drawRadials(gc)
		drawConcentricCircles(gc)
		drawSweep(gc, i)

		for j := range circles {
			drawMovingCircle(gc, circles[j].xpos, circles[j].ypos, circles[j].radius, circles[j].brightness)
		}
		updateCircles(i)

		drawMaskingCircle(gc)

		gl.Flush() /* single buffered, so needs a flush. */
		time.Sleep(2 * time.Millisecond)
		w.SwapBuffers()
		//time.Sleep(200 * time.Millisecond)
	}
}

// Draws radial lines in radar.
func drawRadials(gc *draw2dgl.GraphicContext) {
	gc.Save()
	gc.Translate(width/2, height/2)
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})
	for i := 0.0; i < 360; i = i + 45 { // go from 0 to 360 degrees in 45 degree steps
		gc.Save()                        // keep rotations temporary
		gc.Rotate(i * (math.Pi / 180.0)) // rotate by degrees from i
		gc.MoveTo(0, 0)
		// uses the smaller of width/height to determine how long to draw the radials
		if width > height {
			gc.LineTo(height*0.5, 0)
		} else {
			gc.LineTo(width*0.5, 0)
		}
		gc.Stroke()
		gc.Restore()
	}
	gc.Restore()
}

// Draws concentric circles in radar.
func drawConcentricCircles(gc *draw2dgl.GraphicContext) {
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})

	var radius float64
	// uses the smaller of width/height to determine how big to make the circles
	if width > height {
		radius = height / 2
	} else {
		radius = width / 2
	}
	x := width / 2
	y := height / 2
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0) //clockwise in radians
	gc.SetLineCap(draw2d.ButtCap)
	for i := 1.0; i > 0; i = i - 0.3 { // reduction factor for concentric circles
		gc.MoveTo(x+math.Cos(startAngle)*radius, y+math.Sin(startAngle)*radius)
		gc.ArcTo(x, y, radius*i, radius*i, startAngle, sweepAngle)
		gc.Stroke()
	}
}

// Draws the radar sweep.
func drawSweep(gc *draw2dgl.GraphicContext, angle int) {
	var length float64
	// uses the smaller of width/height to determine how long to draw the sweep radials
	if width > height {
		length = height / 2
	} else {
		length = width / 2
	}
	// draw multiple radials to create fading sweep
	for i := 0; i < 60; i++ {
		gc.Save()
		gc.Translate(width/2, height/2)
		gc.SetLineWidth(1)
		gc.SetStrokeColor(color.RGBA{0, uint8(250 - 3*i), 0, 0xff}) // draws each radial less bright
		gc.Rotate((float64(angle) - float64(i)*0.5) * (math.Pi / 180.0))
		gc.MoveTo(0, 0)
		gc.LineTo(length, 0)
		gc.Stroke()
		gc.Restore()
	}
}

// Draws the moving circles.
//	x: x coordinate of center of circle
//	y: y coordinate of center of circle
//	radius: radius of circle
//	brightness: brightness of circle
func drawMovingCircle(gc *draw2dgl.GraphicContext, x, y, radius float64, brightness int) {
	gc.SetLineWidth(4)
	gc.SetStrokeColor(color.RGBA{0, uint8(brightness), 0, 0xff})
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0) //clockwise in radians
	gc.SetLineCap(draw2d.ButtCap)
	gc.MoveTo(x+math.Cos(startAngle)*radius, y+math.Sin(startAngle)*radius)
	gc.ArcTo(x, y, radius, radius, startAngle, sweepAngle)
	gc.Stroke()
}

// Draws a black circle to mask the moving circles that move outside of the outer radar circle
func drawMaskingCircle(gc *draw2dgl.GraphicContext) {
	var radius float64
	// uses the smaller of width/height to determine how big to make the circle
	if width > height {
		radius = height / 2
	} else {
		radius = width / 2
	}
	x := width / 2
	y := height / 2
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0) //clockwise in radians
	gc.SetLineWidth(radius)
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
	gc.SetLineCap(draw2d.ButtCap)
	gc.MoveTo(x+math.Cos(startAngle)*radius, y+math.Sin(startAngle)*radius)
	gc.ArcTo(x, y, radius+(radius/2)+1, radius+(radius/2)+1, startAngle, sweepAngle)
	gc.Stroke()
}
