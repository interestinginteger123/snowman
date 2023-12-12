package main

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

func randomPosition(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func main() {
	a := app.App()
	scene := core.NewNode()

	gui.Manager().Set(scene)

	cam := camera.New(1)
	cam.SetPosition(0, -1, 5)
	scene.Add(cam)

	camera.NewOrbitControl(cam)

	onResize := func(evname string, ev interface{}) {
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	//Snowman Body
	bodyGeom := geometry.NewSphere(1, 32, 32)
	bodyMat := material.NewStandard(math32.NewColor("White"))
	bodyMesh := graphic.NewMesh(bodyGeom, bodyMat)
	scene.Add(bodyMesh)

	//Left Arm
	leftArmGeom := geometry.NewCylinder(0.07, 0.7, 16, 16, true, false) // Narrower and longer arm
	leftArmMat := material.NewStandard(math32.NewColor("Brown"))
	leftArmMesh := graphic.NewMesh(leftArmGeom, leftArmMat)
	leftArmMesh.SetPosition(-0.8, 0.5, 0) // Lower the arm slightly
	leftArmMesh.SetScale(1, 1.7, 1)       // Make it longer
	leftArmMesh.SetRotationZ(math32.Pi / 3.75)
	scene.Add(leftArmMesh)

	//Left Fingers
	for i := 0; i < 3; i++ {
		fingerGeom := geometry.NewCylinder(0.03, 0.3, 16, 16, true, false) // Adjust finger size
		fingerMat := material.NewStandard(math32.NewColor("Brown"))
		fingerMesh := graphic.NewMesh(fingerGeom, fingerMat)

		// Calculate the position of each finger at the end of the left arm
		fingerX := float32(-1.20)                         // Adjust the initial X position for each finger
		fingerY := float32(0.90)                          // Set all fingers at the same height
		fingerZ := 0.08 - float32(i)*0.2 + float32(i)*0.1 // Adjust the spacing between fingers

		fingerMesh.SetPosition(fingerX, fingerY, fingerZ)

		if i == 0 {
			fingerMesh.SetRotationX(math32.Pi / 6) // Rotate by 30 degrees (Pi/6 radians)
		}

		if i == 2 {
			fingerMesh.SetRotationX(-math32.Pi / 6) // Rotate by -30 degrees (-Pi/6 radians)
		}
		armRotationZ := leftArmMesh.Rotation().Z
		fingerMesh.SetRotationZ(armRotationZ)

		scene.Add(fingerMesh)
	}

	//Right arm
	rightArmGeom := geometry.NewCylinder(0.07, 0.7, 16, 16, true, false) // Narrower and longer arm
	rightArmMat := material.NewStandard(math32.NewColor("Brown"))
	rightArmMesh := graphic.NewMesh(rightArmGeom, rightArmMat)
	rightArmMesh.SetPosition(0.8, 0.5, 0) // Lower the arm slightly
	rightArmMesh.SetScale(1, 1.7, 1)      // Make it longer
	rightArmMesh.SetRotationZ(-math32.Pi / 3.75)
	scene.Add(rightArmMesh)

	//Right Fingers
	for i := 0; i < 3; i++ {
		fingerGeom := geometry.NewCylinder(0.03, 0.3, 16, 16, true, false) // Adjust finger size
		fingerMat := material.NewStandard(math32.NewColor("Brown"))
		fingerMesh := graphic.NewMesh(fingerGeom, fingerMat)
		fingerX := float32(1.22)
		fingerY := float32(1.0) - 0.05
		fingerZ := 0.10 - float32(i)*0.2 + float32(i)*0.1
		fingerMesh.SetPosition(fingerX, fingerY, fingerZ)

		if i == 0 {
			fingerMesh.SetRotationX(math32.Pi / 6)
		}

		if i == 2 {
			fingerMesh.SetRotationX(-math32.Pi / 6)
		}

		armRotationZ := rightArmMesh.Rotation().Z
		fingerMesh.SetRotationZ(armRotationZ)

		scene.Add(fingerMesh)
	}

	//Snowman Head
	headGeom := geometry.NewSphere(0.8, 32, 32)
	headMat := material.NewStandard(math32.NewColor("White"))
	headMesh := graphic.NewMesh(headGeom, headMat)
	headMesh.SetPosition(0, 1.5, 0)
	scene.Add(headMesh)

	// Scarf
	scarfGeom := geometry.NewCylinder(0.65, 0.2, 16, 16, true, false) // Adjust cylinder parameters for the larger scarf
	scarfMat := material.NewStandard(math32.NewColor("Red"))          // You can use any color for the scarf
	scarfMesh := graphic.NewMesh(scarfGeom, scarfMat)
	scarfMesh.SetPosition(0.0, 0.8, 0)
	scene.Add(scarfMesh)

	// Right end of the scarf
	rightEndGeom := geometry.NewTorus(0.25, 0.08, 16, 32, math32.Pi)
	rightEndMat := material.NewStandard(math32.NewColor("Red"))
	rightEndMesh := graphic.NewMesh(rightEndGeom, rightEndMat)
	rightEndMesh.SetPosition(0.3, 0.7, 0.5)  // Adjust position closer and to the front
	rightEndMesh.SetRotationX(math32.Pi / 6) // Adjust the rotation for a downward angle
	scene.Add(rightEndMesh)

	// Loose end
	looseEndGeom := geometry.NewBox(0.1, 0.1, 0.5) // Adjust the size of the rectangle
	looseEndMat := material.NewStandard(math32.NewColor("Red"))
	looseEndMesh := graphic.NewMesh(looseEndGeom, looseEndMat)
	looseEndMesh.SetPosition(0.55, 0.50, 0.55) // Adjust the position
	looseEndMesh.SetRotationX(math32.Pi / 2)   // Rotate the rectangle to align properly
	scene.Add(looseEndMesh)

	//Snowman Bottom
	bottomGeom := geometry.NewSphere(1.2, 32, 32)
	bottomMat := material.NewStandard(math32.NewColor("White"))
	bottomMesh := graphic.NewMesh(bottomGeom, bottomMat)
	bottomMesh.SetPosition(0, -1.5, 0)
	scene.Add(bottomMesh)

	//Eyes
	eyeGeom := geometry.NewSphere(0.15, 16, 16)
	eyeMat := material.NewStandard(math32.NewColor("Black"))
	leftEye := graphic.NewMesh(eyeGeom, eyeMat)
	rightEye := graphic.NewMesh(eyeGeom, eyeMat)
	leftEye.SetPosition(-0.4, 1.8, 0.5)
	rightEye.SetPosition(0.4, 1.8, 0.5)
	scene.Add(leftEye)
	scene.Add(rightEye)

	//Nose
	noseGeom := geometry.NewCone(0.15, 0.6, 16, 16, true)
	noseMat := material.NewStandard(math32.NewColor("Orange"))
	noseMesh := graphic.NewMesh(noseGeom, noseMat)
	noseMesh.SetRotationX(math32.Pi / 2)
	noseMesh.SetPosition(0, 1.5, 0.9)
	scene.Add(noseMesh)

	//Mouth
	mouthGeom := geometry.NewTorus(0.2, 0.05, 16, 32, math32.Pi)
	mouthMat := material.NewStandard(math32.NewColor("Black"))
	mouthMesh := graphic.NewMesh(mouthGeom, mouthMat)
	mouthMesh.SetRotationZ(math32.Pi)
	mouthMesh.SetPosition(0, 1.2, 0.7)
	scene.Add(mouthMesh)

	//Buttons
	buttonGeom := geometry.NewSphere(0.1, 16, 16)
	buttonMat := material.NewStandard(math32.NewColor("Black"))
	button1 := graphic.NewMesh(buttonGeom, buttonMat)
	button2 := graphic.NewMesh(buttonGeom, buttonMat)
	button3 := graphic.NewMesh(buttonGeom, buttonMat)
	button1.SetPosition(0, 0.5, 0.9)
	button2.SetPosition(0, 0, 1)
	button3.SetPosition(0, -0.5, 0.9)
	scene.Add(button1)
	scene.Add(button2)
	scene.Add(button3)

	//Top hat
	hatBaseGeom := geometry.NewCylinder(0.4, 0.6, 32, 34, true, false)
	hatBaseMat := material.NewStandard(math32.NewColor("Black"))
	hatBaseMesh := graphic.NewMesh(hatBaseGeom, hatBaseMat)
	hatBaseMesh.SetPosition(0, 2.5, 0)
	scene.Add(hatBaseMesh)

	//More Top hat
	hatCircleGeom := geometry.NewCylinder(0.6, 0.01, 32, 34, true, false)
	hatCircleMat := material.NewStandard(math32.NewColor("Black"))
	hatCircleMesh := graphic.NewMesh(hatCircleGeom, hatCircleMat)
	hatCircleMesh.SetPosition(0, 2.5-0.3, 0)
	scene.Add(hatCircleMesh)

	//Lights camera action baby
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	//Snow
	numSnowflakes := 1000
	snowflakes := make([]*graphic.Mesh, numSnowflakes)

	for i := 0; i < numSnowflakes; i++ {
		snowflakeGeom := geometry.NewSphere(0.01, 8, 8)
		snowflakeMat := material.NewStandard(math32.NewColor("White"))
		snowflakeMesh := graphic.NewMesh(snowflakeGeom, snowflakeMat)

		x := rand.Float32()*10 - 5
		y := rand.Float32()*5 + 5
		z := rand.Float32()*10 - 5

		snowflakeMesh.SetPosition(x, y, z)
		scene.Add(snowflakeMesh)
		snowflakes[i] = snowflakeMesh
	}

	//floor
	floorGeom := geometry.NewPlane(10, 10)
	floorMat := material.NewStandard(math32.NewColor("White"))
	floorMesh := graphic.NewMesh(floorGeom, floorMat)
	floorMesh.SetRotationX(-math32.Pi / 2)
	floorMesh.SetPosition(0, -2, 0)
	scene.Add(floorMesh)

	//Merry Christmas
	label := gui.NewLabel("Merry Christmas Ello Group from Shaun")
	label.SetFontSize(24)
	label.SetPosition(10, 10)
	scene.Add(label)

	scene.Add(helper.NewAxes(0.5))
	a.Gls().ClearColor(0.0, 0.0, 0.2, 1.0)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		renderer.Render(scene, cam)
		for _, snowflake := range snowflakes {
			speed := 1
			newY := snowflake.Position().Y - float32(speed)*float32(deltaTime.Seconds())
			snowflake.SetPositionY(newY)

			if newY < -1 {
				snowflake.SetPositionY(5)
			}
		}

		if err := renderer.Render(scene, cam); err != nil {
			fmt.Println("Error rendering scene:", err)
		}
	})
}
