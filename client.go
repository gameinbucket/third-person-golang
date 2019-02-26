package main

import (
	g "./graphics"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
	"runtime"
)

const (
	Key1     = 0
	Key2     = 1
	Key3     = 2
	KeyCount = 3

	perspectiveNear = float32(0.01)
	perspectiveFar  = float32(256.0)
)

var (
	canvasWidth        = 11 * 75
	canvasHeight       = 7 * 75
	window             *glfw.Window
	aspectRatio        float32
	perspectiveFov     = float32(math.Tan(float64(60.0) * math.Pi / 360.0))
	keyMap             [KeyCount]bool
	programPreTexture2 uint32
	programPreTexture3 uint32
	programPreModel    uint32
	programPreBlocks   uint32
	myOrthographic     = make([]float32, 16)
	myPerspective      = make([]float32, 16)
	genericBuffer      *g.Buffer
	modelBuffer        [2]*g.Buffer
	gpuIndex           = 0
	imgAtlas           uint32
	imgAtlasArray      uint32
	lense              *Camera
	isWireFrame        = false
	isCullEnabled      = false
	toggleMode         = false
	you                = new(Entity)
)

func Resize(w *glfw.Window, width, height int) {
	canvasWidth = width
	canvasHeight = height
	aspectRatio = float32(canvasWidth) / float32(canvasHeight)

	g.Orthographic(myOrthographic, 0.0, 0.0, float32(canvasWidth), float32(canvasHeight), 0, 1)
	g.Perspective(myPerspective, perspectiveFov, perspectiveNear, perspectiveFar, aspectRatio)
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	/* window */

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	local, err := glfw.CreateWindow(canvasWidth, canvasHeight, "Sword and Sigil", nil, nil)
	if err != nil {
		panic(err)
	}
	window = local

	videoMode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((videoMode.Width-canvasWidth)>>1, (videoMode.Height-canvasHeight)>>1)

	window.MakeContextCurrent()
	window.SetSizeCallback(Resize)
	glfw.SwapInterval(1)

	/* gl */

	if err := gl.Init(); err != nil {
		panic(err)
	}
	g.SetClearColor(0.0, 0.0, 0.0)
	gl.Disable(gl.BLEND)
	gl.Disable(gl.CULL_FACE)

	/* graphics */

	glerr := gl.GetError()
	if glerr != gl.NO_ERROR {
		fmt.Printf("gl error 0x%x\n", glerr)
	}
	programPreTexture2 = g.MakeProgram("pre-texture2", false)
	programPreTexture3 = g.MakeProgram("pre-texture3", false)
	programPreModel = g.MakeProgram("pre-model", false)
	programPreBlocks = g.MakeProgram("pre-blocks", false)
	imgAtlas = g.MakeTexture("atlas.png", true, false)
	imgAtlasArray = g.MakeTextureArray([]string{"grass.png", "stone.png"}, 16, 16, false, false)
	genericBuffer = g.NewBuffer(
		[]g.Attribute{
			{Typ: gl.FLOAT, Size: 2},
			{Typ: gl.FLOAT, Size: 2}}, 4, 6)
	modelBuffer[0] = g.NewBuffer(
		[]g.Attribute{
			{Typ: gl.FLOAT, Size: 3},
			{Typ: gl.FLOAT, Size: 3},
			{Typ: gl.FLOAT, Size: 2}}, 4*100, 6*100)
	modelBuffer[1] = g.NewBuffer(
		[]g.Attribute{
			{Typ: gl.FLOAT, Size: 3},
			{Typ: gl.FLOAT, Size: 3},
			{Typ: gl.FLOAT, Size: 2}}, 4*100, 6*100)

	/* world */

	you.run = float32(0.1)
	you.Y = 2.5
	you.Type = EntityHuman
	you.mesh = new(g.Model)
	makeBiped(you.mesh)

	lense = new(Camera)
	lense.radius = 5.0
	lense.x = 2.0
	lense.y = 3.0
	lense.z = 2.0
	lense.ry = 1.07
	lense.follow = you

	MakeWorld()
	BuildWorld()

	/* run */

	Resize(window, canvasWidth, canvasHeight)
	window.Show()
	for !window.ShouldClose() {
		glfw.PollEvents()
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}
		Update()
		Render()
		window.SwapBuffers()
	}
}

func Update() {
	turn := float32(0.05)
	if window.GetKey(glfw.KeyJ) == glfw.Press {
		lense.ry -= turn
		if lense.ry < 0 {
			lense.ry += math.Pi * 2
		}
	}
	if window.GetKey(glfw.KeyL) == glfw.Press {
		lense.ry += turn
		if lense.ry >= math.Pi*2 {
			lense.ry -= math.Pi * 2
		}
	}
	if window.GetKey(glfw.KeyI) == glfw.Press {
		lense.rx += turn
		if lense.rx >= math.Pi*2 {
			lense.rx -= math.Pi * 2
		}
	}
	if window.GetKey(glfw.KeyK) == glfw.Press {
		lense.rx -= turn
		if lense.rx < 0 {
			lense.rx += math.Pi * 2
		}
	}

	if window.GetKey(glfw.Key1) == glfw.Press {
		if keyMap[Key1] == false {
			if isWireFrame {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			}
			isWireFrame = !isWireFrame
		}
		keyMap[Key1] = true
	} else {
		keyMap[Key1] = false
	}
	if window.GetKey(glfw.Key2) == glfw.Press {
		if keyMap[Key2] == false {
			if isCullEnabled {
				gl.Disable(gl.CULL_FACE)
			} else {
				gl.Enable(gl.CULL_FACE)
			}
			isCullEnabled = !isCullEnabled
		}
		keyMap[Key2] = true
	} else {
		keyMap[Key2] = false
	}
	if window.GetKey(glfw.Key3) == glfw.Press {
		if keyMap[Key3] == false {
			toggleMode = !toggleMode
		}
		keyMap[Key3] = true
	} else {
		keyMap[Key3] = false
	}

	you.Update()
	lense.Update()
}

func Render() {
	gpuIndex = 1 - gpuIndex
	mb := modelBuffer[gpuIndex]
	g.ClearColorDepth()

	g.SetFramebuffer(0)
	g.SetView(0, 0, int32(canvasWidth), int32(canvasHeight))

	gl.Enable(gl.DEPTH_TEST)
	g.SetProgram(programPreBlocks)
	g.SetPerspective(myPerspective, -lense.x, -lense.y, -lense.z, lense.rx, lense.ry)
	g.SetMVP()
	g.SetTextureArray0(imgAtlasArray)
	gx := int(lense.x) >> chunkShift
	gy := int(lense.y) >> chunkShift
	gz := int(lense.z) >> chunkShift
	ReadyFrustum()
	Occlude(gx, gy, gz)
	RenderWorld(gx, gy, gz)

	g.SetProgram(programPreModel)
	g.SetMVP()
	mb.Zero()
	you.Render(mb)
	g.SetTextureArray0(imgAtlas)
	mb.Draw()

	gl.Disable(gl.DEPTH_TEST)
	g.SetProgram(programPreTexture2)
	g.SetOrthographic(myOrthographic, 0.0, 0.0)
	g.SetMVP()
	genericBuffer.Zero()
	genericBuffer.RenderImage(0.0, 0.0, 64.0, 64.0, 0.0, 0.0, 1.0, 1.0)
	g.SetTexture0(imgAtlas)
	genericBuffer.Draw()
}
