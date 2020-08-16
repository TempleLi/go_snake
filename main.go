package main

import (
	"fmt"
	"games/snake"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// warning!
// install opengl dependency
// ref to: https://github.com/faiface/pixel#requirements

func main() {
	pixelgl.Run(run)
}
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "贪吃蛇",
		Bounds: pixel.R(0, 0, 1000, 1000),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	s:=snake.NewSnake(5,50,40)
	imd := imdraw.New(nil)
	s.Restart()
	s.Run()
	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Color = colornames.Red
		if win.JustPressed(pixelgl.KeyLeft) {
			s.SetDirection(snake.Left)
		}
		if win.JustPressed(pixelgl.KeyRight) {
			s.SetDirection(snake.Right)
		}
		if win.JustPressed(pixelgl.KeyDown) {
			s.SetDirection(snake.Bottom)
		}
		if win.JustPressed(pixelgl.KeyUp) {
			s.SetDirection(snake.Top)
		}
		if win.JustPressed(pixelgl.KeySpace){
			s.Pause()
		}
		if win.JustPressed(pixelgl.KeyEnter){
			s.Restart()
		}
		// desc game state

		// draw header
		if s.State()==snake.Over{
			imd.Color=colornames.Red
		}else{
			imd.Color=colornames.Black
		}
		imd.Push(pixel.V(0,800))
		imd.Push(pixel.V(1000,800))
		imd.Push(pixel.V(1000,1000))
		imd.Push(pixel.V(0,1000))
		imd.Polygon(0)

		imd.Color=colornames.Gray
		imd.Push(pixel.V(0,0))
		imd.Push(pixel.V(0,800))
		imd.Push(pixel.V(1000,800))
		imd.Push(pixel.V(1000,0))
		imd.Polygon(0)
		imd.Draw(win)
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 880), basicAtlas)
		if s.State()==snake.Over{
			basicTxt.Color=colornames.Black
		}else{
			basicTxt.Color=colornames.Darkorange
		}
		var state string
		switch s.State(){
		case snake.Idle:
			state="Idle"
		case snake.Running:
			state="Running"
		case snake.Pause:
			state="Pause"
		case snake.Over:
			state="Game Over"
		}
		_,_=fmt.Fprintf(basicTxt, "%s",state)
		basicTxt.Draw(win,  pixel.IM.Scaled(basicTxt.Orig, 4))
		basicTxt = text.New(pixel.V(400, 840), basicAtlas)
		basicTxt.Color=colornames.Beige
		_,_=fmt.Fprintf(basicTxt, "Score:%3d  Steps:%4d",s.Score(),s.Steps())
		basicTxt.Draw(win,  pixel.IM.Scaled(basicTxt.Orig, 3))

		imd=s.GetDraw(20)
		imd.Draw(win)
		win.Update()
	}
}
