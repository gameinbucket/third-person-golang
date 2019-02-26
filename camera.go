package main

import (
	"math"
)

type Camera struct {
	follow *Entity
	rx     float32
	ry     float32
	radius float32
	x      float32
	y      float32
	z      float32
}

func (o *Camera) Update() {
	sinX := float32(math.Sin(float64(o.rx)))
	cosX := float32(math.Cos(float64(o.rx)))
	sinY := float32(math.Sin(float64(o.ry)))
	cosY := float32(math.Cos(float64(o.ry)))
	x := o.follow.X - o.radius*cosX*sinY
	y := o.follow.Y + o.radius*sinX
	z := o.follow.Z + o.radius*cosX*cosY
	o.x, o.y, o.z = raycastXYZ(o.follow.X, o.follow.Y, o.follow.Z, x, y, z)
}
