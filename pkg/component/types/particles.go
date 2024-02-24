package types

import (
	"encoding/xml"
	c "github.com/6ixisgood/matrix-ticker/pkg/component/common"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type Particle struct {
	X, Y   float64
	SpeedX float64
	SpeedY float64
	Color  color.RGBA
	Radius float64
}

type GravityParticle struct {
	X, Y  float64
	Force float64
	Color color.RGBA
}

type GravityParticles struct {
	c.BaseComponent

	XMLName       xml.Name `xml:"gravity-particles"`
	Particles     []Particle
	GravityPoints []GravityParticle
}

func (gp *GravityParticles) Init() {
	gp.BaseComponent.Init()

	// Sample initialization for particles and gravity points
	// In a real-world scenario, these could be populated based on the music or other parameters
	for i := 0; i < 50; i++ {
		gp.Particles = append(gp.Particles, Particle{
			X:      float64(rand.Intn(gp.Width())),
			Y:      float64(rand.Intn(gp.Height())),
			SpeedX: float64(rand.Intn(5) - 2), // Random speed between -2 and 2
			SpeedY: float64(rand.Intn(5) - 2),
			Color:  color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
			Radius: 2,
		})
	}

	gp.GravityPoints = append(gp.GravityPoints, GravityParticle{
		X:     float64(gp.Width()) / 2,
		Y:     float64(gp.Height()) / 2,
		Force: 5,
		Color: color.RGBA{255, 0, 0, 255},
	})
}

func (gp *GravityParticles) Render() image.Image {
	gp.Ctx.SetColor(color.RGBA{0, 0, 0, 255})
	gp.Ctx.Clear()

	// Update and render particles
	for i, particle := range gp.Particles {
		for _, g := range gp.GravityPoints {
			// Calculate the force exerted by this gravity point
			distX := g.X - particle.X
			distY := g.Y - particle.Y
			distance := math.Sqrt(distX*distX + distY*distY)

			// The closer the particle is to the gravity point, the stronger the pull
			force := g.Force / (distance + 1) // +1 to avoid division by zero
			particle.SpeedX += force * distX / distance
			particle.SpeedY += force * distY / distance
		}

		particle.X += particle.SpeedX
		particle.Y += particle.SpeedY

		// Boundary checking
		if particle.X < 0 || particle.X > float64(gp.Width()) {
			particle.SpeedX = -particle.SpeedX
		}
		if particle.Y < 0 || particle.Y > float64(gp.Height()) {
			particle.SpeedY = -particle.SpeedY
		}

		// Update the particle in the slice
		gp.Particles[i] = particle

		// Draw the particle
		gp.Ctx.SetColor(particle.Color)
		gp.Ctx.DrawCircle(particle.X, particle.Y, particle.Radius)
		gp.Ctx.Fill()
	}

	return gp.Ctx.Image()
}

func init() {
	c.RegisterComponent("gravity-particles", func() c.Component { return &GravityParticles{} })
}
