package scene

import (
	"fmt"
	. "github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"sync"
)

type scene interface {
	Init()
	Update(dt float64)
	Draw(*ebiten.Image)
	OnCommand(command string)
	Destroy()
}

type Manager struct {
	mu sync.Mutex

	//name of current scene, change async
	current string

	paused         bool
	pauseSceneName string

	scenes map[string]scene
	inited map[string]bool

	actionQ chan func()
}

func NewManager() *Manager {
	res := &Manager{
		scenes:  make(map[string]scene),
		inited:  make(map[string]bool),
		actionQ: make(chan func(), 3),
		paused:  true, //same as network client state
	}
	go actionRun(res)
	return res
}

//manager goroutine cycle
func actionRun(m *Manager) {
	for f := range m.actionQ {
		m.mu.Lock()
		f()
		m.mu.Unlock()
	}
}

//main cycle
func (m *Manager) UpdateAndDraw(dt float64, image *ebiten.Image, doDraw bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	logStr := fmt.Sprintf("manager.UpdateAndDraw %v paused: %v inited: %v",
		m.current, m.paused, m.inited[m.current])
	defer LogFunc(logStr)()

	if doDraw {
		image.Clear()
	}

	if m.current == "" {
		return
	}

	actualScene := m.current
	if (m.paused || !m.inited[m.current]) && m.pauseSceneName != "" {
		actualScene = m.pauseSceneName
	}

	if !m.inited[actualScene] {
		Log(LVL_ERROR, "trying to update not inited scene", m.current, "while actual ", actualScene)
		return
	}

	scene := m.scenes[actualScene]

	scene.Update(dt)

	if doDraw {

		scene.Draw(image)
	}
}

//Main or Network loop
func (m *Manager) Install(name string, Scene scene, inited bool) {
	m.actionQ <- func() {
		m.install(name, Scene, inited)
	}
}

//Main or Network loop
func (m *Manager) Delete(name string) {
	m.actionQ <- func() {
		m.delete(name)
	}
}

//Main or Network loop
func (m *Manager) Activate(name string, needReInit bool) {
	m.actionQ <- func() {
		m.activate(name, needReInit)
	}
}

//Main or Network loop
func (m *Manager) Init(name string) {
	//we do this with non-current scene
	//so we do it NOW
	m.init(name)
}

//Main or Network loop
func (m *Manager) SetAsPauseScene(pauseSceneName string) {
	m.actionQ <- func() {
		m.setAsPauseScene(pauseSceneName)
	}
}

//Main or Network loop
func (m *Manager) SetPaused(paused bool) {
	m.actionQ <- func() {
		m.setPaused(paused)
	}
}

//Main or Network loop
func (m *Manager) OnCommand(command string) {
	m.actionQ <- func() {
		m.onCommand(command)
	}
}

//Main or Network loop
func (m *Manager) WaitDone() {
	done := make(chan struct{})
	m.actionQ <- func() {
		close(done)
	}
	<-done
}

//manager goroutine cycle
func (m *Manager) install(name string, Scene scene, inited bool) {
	if prev, ok := m.scenes[name]; ok {
		prev.Destroy()
	}

	m.scenes[name] = Scene
	m.inited[name] = inited
}

//manager goroutine cycle
func (m *Manager) delete(name string) {
	if prev, ok := m.scenes[name]; ok {
		prev.Destroy()
		delete(m.scenes, name)

		if m.current == name {
			m.current = ""
		}
	}
}

//manager goroutine cycle
func (m *Manager) activate(name string, needReInit bool) {
	if _, ok := m.scenes[name]; !ok {
		panic("can't activate scene " + name)
	}

	m.current = name
	if needReInit {
		m.inited[name] = false
	}
}

//manager goroutine cycle
func (m *Manager) init(name string) {
	m.mu.Lock()
	scene, ok := m.scenes[name]
	m.mu.Unlock()

	if ok {
		scene.Init()

		m.mu.Lock()
		m.inited[name] = true
		m.mu.Unlock()
	}
}

//manager goroutine cycle
func (m *Manager) onCommand(command string) {
	scene, ok := m.scenes[m.current]

	if ok {
		scene.OnCommand(command)
	}
}

//manager goroutine cycle
func (m *Manager) setAsPauseScene(pauseSceneName string) {
	m.pauseSceneName = pauseSceneName
}

//manager goroutine cycle
func (m *Manager) setPaused(paused bool) {
	m.paused = paused
}
