package snake

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
	"sync"
	"time"
)
type Direction uint8
const (
	Left = iota
	Top
	Right
	Bottom
)

const (
	Idle =iota
	Running
	Pause
	Over
)

type Snake struct{
	initSize int
	steps int
	grows int
	body [][2]int

	width     int
	height    int
	state     int
	direction Direction
	lock      sync.Mutex
	profit    map[[2]int]int
}
func NewSnake(size int,width int,height int)*Snake{
	if width<=0 || height<=0 || size<=0{
		panic("illegal argument")
	}
	s:=Snake{
		width:     width,
		height:    height,
		direction: Right,
		initSize:  size,
		state:     Idle,
	}
	s.Reset()
	return &s
}
func (s *Snake)Pause(){
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.state==Running{
		s.state=Pause
	}else if s.state==Pause{
		s.state=Running
	}
}
func (s *Snake)Restart(){
	s.lock.Lock()
	defer s.lock.Unlock()
	s.reset()
	s.state=Running
}

// please call this only once
func(s *Snake) Run(){
	go func(){
		ticker:=time.NewTicker(time.Millisecond*100)
		for {
			<-ticker.C
			s.lock.Lock()
			s.step()
			s.lock.Unlock()
		}
	}()
}

func(s *Snake) Time(){
}

func(s *Snake) Score()int{
	s.lock.Lock()
	s.lock.Unlock()
	return len(s.body)
}
func (s *Snake)step(){
	if s.state!=Running{
		return
	}
	s.steps++
	if s.grows>0{
		s.grows--
		end:=s.newEnd()
		s.body=append(s.body,end)
	}else{
		for i:=0;i<len(s.body)-1;i++{
			s.body[i][0]=s.body[i+1][0]
			s.body[i][1]=s.body[i+1][1]
		}
		end:=s.newEnd()
		s.body[len(s.body)-1]=end
	}
	end:=s.body[len(s.body)-1]
	if s.profit[end]>0{
		s.grows=s.profit[end]
		delete(s.profit,end)
	}
	if len(s.profit)<4{
		s.randProfit(12)
	}

	if s.isDead(){
		s.state= Over
	}
}
func (s *Snake) newEnd()[2]int{
	offsetX:=0
	offsetY:=0
	switch s.direction{
	case Top:
		offsetY=1
	case Right:
		offsetX=1
	case Bottom:
		offsetY=-1
	case Left:
		offsetX=-1
	}
	x:=s.body[len(s.body)-1][0]+offsetX
	y:=s.body[len(s.body)-1][1]+offsetY
	return [2]int{x,y}
}

func (s *Snake) isDead()bool{
	size:=len(s.body)
	lastP:=s.body[size-1]
	if lastP[0]>=s.width || lastP[0]<0{
		return true
	}
	if lastP[1]>=s.height || lastP[1]<0{
		return true
	}
	for i:=0;i<size-1;i++{
		if lastP==s.body[i]{
			return true
		}
	}
	return false
}

func (s *Snake) GetDraw(span int)*imdraw.IMDraw{
	s.lock.Lock()
	defer s.lock.Unlock()

	imd := imdraw.New(nil)
	for i,p:=range s.body{
		if i==len(s.body)-1{
			drawPoint(imd,p,span,colornames.Black)
		}else{
			drawPoint(imd,p,span,colornames.Red)
		}
	}
	for p,_:=range s.profit{
		drawPoint(imd,p,span,colornames.Orange)
	}
	return imd
}
func drawPoint(draw *imdraw.IMDraw,p[2]int,spanInt int,_color color.Color){
	draw.Color = _color
	span:=float64(spanInt)
	x:=float64(p[0])*span
	y:=float64(p[1])*span
	draw.Push(pixel.V(x,y))
	draw.Push(pixel.V(x+span,y))
	draw.Push(pixel.V(x+span,y+span))
	draw.Push(pixel.V(x,y+span))
	draw.Polygon(0)

}
func (s *Snake)SetDirection(dir Direction){
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.state!=Running{
		return
	}
	if s.direction==Left && dir==Right{
		return
	}
	if s.direction==Right && dir==Left{
		return
	}
	if s.direction==Top && dir==Bottom{
		return
	}
	if s.direction==Bottom && dir==Top{
		return
	}
	s.direction=dir
}

func (s *Snake) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.reset()
}
func (s *Snake)reset(){
	s.direction=Right
	s.body=nil
	for i:=0;i<s.initSize;i++{
		s.body=append(s.body,[2]int{i,0})
	}
	s.state=Idle
	s.profit=make(map[[2]int]int)
	s.randProfit(10)
}
func (s *Snake)randProfit(n int){
	for i:=0;i<n;i++{
		x:=rand.Intn(s.width)
		y:=rand.Intn(s.height)
		s.profit[[2]int{x,y}]++
	}
}


func (s *Snake) State() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.state
}

func (s *Snake) Steps() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.steps
}

