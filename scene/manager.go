package scene

import (
	"github.com/hajimehoshi/ebiten"
	"sync"
)

type scene interface {
	Init()
	Update(dt float64)
	Draw(image *ebiten.Image)
	Destroy()
}

type Manager struct{
	mu sync.RWMutex

	//name of current scene, change async
	current string

	scenes map[string]scene
}

func NewManager() *Manager{
	return &Manager{
		scenes: make(map[string]scene),
	}
}

func (m *Manager) Install(name string, Scene scene) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if prev,ok:=m.scenes[name];ok{
		prev.Destroy()
	}

	m.scenes[name] = Scene
}

func (m *Manager) Delete(name string){
	m.mu.Lock()
	defer m.mu.Unlock()

	if prev,ok:=m.scenes[name];ok{
		prev.Destroy()
		delete(m.scenes, name)

		if m.current == name{
			m.current = ""
		}
	}
}

func (m *Manager) Activate(name string){
	m.mu.Lock()
	defer m.mu.Unlock()

	if _,ok:=m.scenes[name]; !ok{
		panic("can't activate scene "+name)
	}

	m.current = name
}

func (m *Manager) UpdateAndDraw (dt float64, image *ebiten.Image, doDraw bool) {
	if m.current!="" {
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	scene:=m.scenes[m.current]

	scene.Update(dt)
	if doDraw {
		image.Clear()
		scene.Draw(image)
	}
}