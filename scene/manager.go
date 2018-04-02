package scene

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"sync"
)

type scene interface {
	Init()
	Update(dt float64)
	Draw(image *ebiten.Image)
	Destroy()
}

type Manager struct {
	mu sync.RWMutex

	//name of current scene, change async
	current string

	paused bool
	pauseSceneName string

	scenes map[string]scene
	inited map[string]bool
}

func NewManager() *Manager {
	return &Manager{
		scenes: make(map[string]scene),
		inited: make(map[string]bool),
	}
}

func (m *Manager) Install(name string, Scene scene, inited bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if prev, ok := m.scenes[name]; ok {
		prev.Destroy()
	}

	m.scenes[name] = Scene
	m.inited[name] = inited
}

func (m *Manager) Delete(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if prev, ok := m.scenes[name]; ok {
		prev.Destroy()
		delete(m.scenes, name)

		if m.current == name {
			m.current = ""
		}
	}
}

func (m *Manager) Activate(name string, needReInit bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.scenes[name]; !ok {
		panic("can't activate scene " + name)
	}

	log.Println("Activate scene: ",name)
	m.current = name
	if needReInit {
		m.inited[name] = false
	}
}

func (m *Manager) Init(name string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if scene, ok := m.scenes[name]; ok {
		scene.Init()
		m.inited[name] = true
	}
}

func (m *Manager) UpdateAndDraw(dt float64, image *ebiten.Image, doDraw bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if doDraw {
		image.Clear()
	}

	if m.current == "" {
		return
	}

	actualScene := m.current
	if m.paused && m.pauseSceneName!="" {
		actualScene = m.pauseSceneName
	}

	if !m.inited[actualScene] {
		log.Println("trying to update not inited scene", m.current)
		return
	}

	scene := m.scenes[actualScene]

	scene.Update(dt)
	if doDraw {
		scene.Draw(image)
	}
}

func (m *Manager) SetOnPauseScene(pauseSceneName string) {
	m.mu.Lock()
	m.pauseSceneName = pauseSceneName
	m.mu.Unlock()
}

func (m *Manager) SetPaused(paused bool) {
	m.mu.Lock()
	m.paused = paused
	m.mu.Unlock()
}