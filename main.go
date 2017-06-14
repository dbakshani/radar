// Open an OpenGl window and display a rectangle using a OpenGl GraphicContext
package main

import (
	"image/color"
	"log"
	"math"
	"runtime"
	//"strconv"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	//"github.com/llgcode/draw2d/draw2dkit"
)

var (
	// global rotation
	rotate        int
	//width, height int
	//width, height = 1024, 1024
	width, height = 512, 512
	floatwidth, floatheight = float64(width), float64(height)
	//redraw        = true
	font          draw2d.FontData
)


// called when the window is reshaped/resized
func reshape(window *glfw.Window, w, h int) {
	//gl.ClearColor(1, 1, 1, 1)
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
	//redraw = true
}

// Ask to refresh
//func invalidate() {
	//redraw = true
//}

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
	for !window.ShouldClose() {
		//if redraw {
			drawContents(window)
			//window.SwapBuffers()
			//redraw = false
		//}
		glfw.PollEvents()
		//		time.Sleep(2 * time.Second)
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

// draws the contents of the window
func drawContents(w *glfw.Window) {
	for i := 0.0; i < 360; i = i + 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LineWidth(1)

		gc := draw2dgl.NewGraphicContext(width, height)

		drawRadials(gc, floatwidth, floatheight, floatwidth, floatheight)
		drawCircles(gc, floatwidth, floatheight, floatwidth, floatheight)
		drawSweep(gc, floatwidth, floatheight, floatwidth, floatheight, i)

		gl.Flush() /* single buffered, so needs a flush. */
		w.SwapBuffers()
		//log.Println("angle: " + strconv.FormatFloat(i, 'f', -1, 32))
		time.Sleep(2 * time.Millisecond)
	}
}

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

func drawCircles(gc *draw2dgl.GraphicContext, xc, yc, width, height float64) {
//	gc.Save()
//	gc.Translate(x/2, y/2)
	gc.SetLineWidth(2)
	gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})
//	for i := 0.0; i < 360; i = i + 45 { // go from 0 to 360 degrees in 45 degree steps
//		gc.Save()                        // keep rotations temporary
//		gc.Rotate(i * (math.Pi / 180.0)) // rotate by degrees on stack from 'for'
//		gc.MoveTo(0, 0)
//		if width > height {
//			gc.LineTo(height * 0.60, 0)
//		} else {
//			gc.LineTo(width * 0.60, 0)
//		}
//		gc.Stroke()
//		gc.Restore()
//	}
//	gc.Restore()

	var radius float64
	if width > height {
		radius = height/2
	} else {
		radius = width/2
	}
	xc = width/2
	yc = height/2
	//radiusX, radiusY := width/2, height/2
	//radiusX, radiusY := width/2, height/2
	//startAngle := 0 * (math.Pi / 180.0) /* angles are specified */
	startAngle := 0.0
	sweepAngle := 360 * (math.Pi / 180.0)     /* clockwise in radians           */
//	gc.SetLineWidth(width / 10)
	gc.SetLineCap(draw2d.ButtCap)
//	gc.SetStrokeColor(image.Black)
	for i := 1.0; i > 0; i = i - 0.3 { // reduction factor for concentric circles
		//gc.MoveTo(xc + math.Cos(startAngle)*radiusX, yc + math.Sin(startAngle)*radiusY)
		gc.MoveTo(xc + math.Cos(startAngle) * radius, yc + math.Sin(startAngle) * radius)
		//gc.ArcTo(xc, yc, radiusX, radiusY, startAngle, sweepAngle)
		gc.ArcTo(xc, yc, radius * i, radius * i, startAngle, sweepAngle)
		gc.Stroke()
	}
}

func drawSweep(gc *draw2dgl.GraphicContext, x, y, width, height, angle float64) {
	//for i := 0.0; i < 40; i = i + 0.5 {
	for i := 0; i < 60; i++ {
		gc.Save()
		gc.Translate(x/2, y/2)
		gc.SetLineWidth(1)
		//gc.SetStrokeColor(color.RGBA{0, 250, 0, 0xff})
		//gc.SetStrokeColor(color.RGBA{0, 250, 0, uint8(255 - i)})
		gc.SetStrokeColor(color.RGBA{0, uint8(250 - 3 * i), 0, 0xff})
		//gc.SetLineDash([]float64{3}, 3)
		//gc.Rotate((angle - float64(float64(i) * 0.5)) * (math.Pi / 180.0))
		gc.Rotate((angle - float64(i) * 0.5) * (math.Pi / 180.0))
		gc.MoveTo(0, 0)
		if width > height {
			gc.LineTo(height * 0.5, 0)
		} else {
			gc.LineTo(width * 0.5, 0)
		}
		gc.Stroke()
		gc.Restore()
	}
}

