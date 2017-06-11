// Open an OpenGl window and display a rectangle using a OpenGl GraphicContext
package main

import (
	"image/color"
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	// global rotation
	rotate        int
	//width, height int
	width, height = 1024, 1024
	redraw        = true
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
	redraw = true
}

// Ask to refresh
func invalidate() {
	redraw = true
}

// draws in the window
func display() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.LineWidth(1)
	gc := draw2dgl.NewGraphicContext(width, height)
	gc.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic})

	gc.BeginPath()
	draw2dkit.RoundedRectangle(gc, 200, 200, 600, 600, 100, 100)

	gc.SetFillColor(color.RGBA{0, 200, 0, 0xff})
	gc.Fill()

	gl.Flush() /* Single buffered, so needs a flush. */
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
	//width, height = 800, 800
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
		if redraw {
			display()
			window.SwapBuffers()
			redraw = false
		}
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
