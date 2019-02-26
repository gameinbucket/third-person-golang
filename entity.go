package main

import (
	g "./graphics"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
)

const (
	EntityHuman = int32(1)
	EntityRat   = int32(2)
)

type Entity struct {
	Type      int32
	Radius    float32
	Height    float32
	X         float32
	Y         float32
	Z         float32
	GX        int32
	GY        int32
	GZ        int32
	MinGX     int32
	MinGY     int32
	MinGZ     int32
	MaxGX     int32
	MaxGY     int32
	MaxGZ     int32
	mesh      *g.Model
	run       float32
	rx        float32
	ry        float32
	register1 float32
}

func (o *Entity) Update() {
	switch o.Type {
	case EntityHuman:
		o.UpdateHuman()
	case EntityRat:
		o.UpdateRat()
	}
}

func (o *Entity) UpdateHuman() {
	rate := float32(0.02)
	pi := float32(math.Pi)
	o.register1 += 0.002
	o.mesh.Bones[bipedHead].LocalRx = 0.5
	o.mesh.Bones[bipedLeftArm].LocalRx = o.register1
	o.mesh.Bones[bipedRightForearm].LocalRx = o.register1
	cos := float32(math.Cos(float64(lense.ry)))
	sin := float32(math.Sin(float64(lense.ry)))
	goalA := float32(999.0)
	goalB := float32(999.0)
	if window.GetKey(glfw.KeyW) == glfw.Press {
		o.X += sin * o.run
		o.Z -= cos * o.run
		goalA = -lense.ry
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		o.X -= sin * o.run
		o.Z += cos * o.run
		if goalA == float32(999.0) {
			goalA = -lense.ry + float32(math.Pi)
		} else {
			goalA = float32(999.0)
		}
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		o.X -= cos * o.run
		o.Z -= sin * o.run
		goalB = -lense.ry + float32(math.Pi)*0.5
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		o.X += cos * o.run
		o.Z += sin * o.run
		if goalB == float32(999.0) { // not good
			goalB = -lense.ry - float32(math.Pi)*0.5
		} else {
			goalB = float32(999.0)
		}
	}
	var goal float32
	if goalA == float32(999.0) {
		if goalB == float32(999.0) {
			goal = float32(999.0)
		} else {
			goal = goalB
		}
	} else {
		if goalB == float32(999.0) {
			goal = goalA
		} else {
			goal = (goalA + goalB) * 0.5 // not good
		}
	}
	if goal != float32(999.0) {
		diff := o.ry - goal
		for diff <= pi {
			diff += pi * 2
		}
		for diff > pi {
			diff -= pi * 2
		}
		if diff < 0 {
			if -diff < rate {
				o.ry = goal
			} else {
				o.ry += rate
			}
		} else {
			if diff < rate {
				o.ry = goal
			} else {
				o.ry -= rate
			}
		}
	}
}

func (o *Entity) UpdateRat() {

}

func (o *Entity) Render(gb *g.Buffer) {
	x := int(o.X)
	y := int(o.Y)
	z := int(o.Z)
	cx := x >> chunkShift
	cy := y >> chunkShift
	cz := z >> chunkShift
	bx := x % chunkSize
	by := y % chunkSize
	bz := z % chunkSize
	red, green, blue := unpackRgb(WorldGetColor(cx, cy, cz, bx, by, bz))
	b := &o.mesh.Bones[bipedBody]
	b.LocalRy = o.ry
	b.WorldX = o.X
	b.WorldY = o.Y
	b.WorldZ = o.Z
	b.RecursivePrePass()
	gb.RenderModel(o.mesh, float32(red)/255.0, float32(green)/255.0, float32(blue)/255.0)
}
