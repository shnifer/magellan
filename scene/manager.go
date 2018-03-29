package scene

import (
	"github.com/hajimehoshi/ebiten"
	"sync"
	"errors"
)

type scene interface {
	Init()
	Update(dt float64)
	Draw(image *ebiten.Image)
	Destroy()
}

type Manager struct{
	mu sync.RWMutex

	//name of current scene
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
	Scene.Init()
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

func (m *Manager) Activate(name string) error{
	m.mu.Lock()
	defer m.mu.Unlock()

	if _,ok:=m.scenes[name]; !ok{
		return errors.New("can't activate scene "+name)
	}

	m.current = name
	return nil
}

func (m *Manager) Update (dt float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current!=""{
		m.scenes[m.current].Update(dt)
	}
}

func (m *Manager) Draw (image *ebiten.Image) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.current!=""{
		m.scenes[m.current].Draw(image)
	}
}