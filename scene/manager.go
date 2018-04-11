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
	}
	go actionRun(res)
	return res
}

func actionRun(m *Manager) {
	for f := range m.actionQ {
		f()
	}
}

func (m *Manager) UpdateAndDraw(dt float64, image *ebiten.Image, doDraw bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
		log.Println("trying to update not inited scene", m.current, "while actual ", actualScene)
		return
	}

	scene := m.scenes[actualScene]

	scene.Update(dt)
	if doDraw {
		scene.Draw(image)
	}
}

func (m *Manager) Install(name string, Scene scene, inited bool) {
	m.actionQ <- func() {
		m.install(name, Scene, inited)
	}
}

func (m *Manager) Delete(name string) {
	m.actionQ <- func() {
		m.delete(name)
	}
}

func (m *Manager) Activate(name string, needReInit bool) {
	m.actionQ <- func() {
		m.activate(name, needReInit)
	}
}

//Init is waiting for init to finish!
func (m *Manager) Init(name string) (done chan struct{}) {
	done = make(chan struct{})
	m.actionQ <- func() {
		m.init(name)
		close(done)
	}
	return done
}

func (m *Manager) SetAsPauseScene(pauseSceneName string) {
	m.actionQ <- func() {
		m.setAsPauseScene(pauseSceneName)
	}
}

func (m *Manager) SetPaused(paused bool) {
	m.actionQ <- func() {
		m.setPaused(paused)
	}
}

func (m *Manager) OnCommand(command string) {
	m.actionQ <- func() {
		m.onCommand(command)
	}
}

func (m *Manager) install(name string, Scene scene, inited bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if prev, ok := m.scenes[name]; ok {
		prev.Destroy()
	}

	m.scenes[name] = Scene
	m.inited[name] = inited
}

func (m *Manager) delete(name string) {
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

func (m *Manager) activate(name string, needReInit bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.scenes[name]; !ok {
		panic("can't activate scene " + name)
	}

	m.current = name
	if needReInit {
		m.inited[name] = false
	}
}

func (m *Manager) init(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if scene, ok := m.scenes[name]; ok {
		scene.Init()
		m.inited[name] = true
	}
}

func (m *Manager) onCommand(command string) {
	m.mu.Lock()
	scene, ok := m.scenes[m.current]
	m.mu.Unlock()

	if ok {
		scene.OnCommand(command)
	}
}

func (m *Manager) setAsPauseScene(pauseSceneName string) {
	m.mu.Lock()
	m.pauseSceneName = pauseSceneName
	m.mu.Unlock()
}

func (m *Manager) setPaused(paused bool) {
	m.mu.Lock()
	m.paused = paused
	m.mu.Unlock()
}
