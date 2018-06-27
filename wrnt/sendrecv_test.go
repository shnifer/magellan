package wrnt

import (
	"testing"
)

type frame struct {
	t    *testing.T
	send *Send
	recv *Recv

	received []string
}

func TestNonConfirmed(t *testing.T) {
	f := newFrame(t)
	f.addAndSend("a", "b")
	f.addAndSend("c")
	f.checkRecv()
}

func TestConfirmed(t *testing.T) {
	f := newFrame(t)
	f.addAndSend("a", "b")
	f.confirm()
	f.addAndSend("c")
	f.checkRecv("a", "b", "c")
}

func TestSenderDrop(t *testing.T) {
	f := newFrame(t)
	f.confirm()
	f.addAndSend("a", "b")
	f.checkRecv("a", "b")
	f.dropSend()
	f.addAndSend("c")
	f.checkRecv("a", "b")
	f.confirm()
	f.addAndSend()
	f.checkRecv("a", "b", "c")
	f.dropSend()
	f.confirm()
	f.checkRecv("a", "b", "c")
	f.addAndSend("d")
	f.checkRecv("a", "b", "c", "d")
}

func TestRecvDrop(t *testing.T) {
	f := newFrame(t)
	f.confirm()
	f.addAndSend("a", "b")
	f.checkRecv("a", "b")
	f.dropRecv()
	f.addAndSend("c")
	f.checkRecv("a", "b", "c")
	f.confirm()
	f.dropRecv()
	f.addAndSend("d", "e")
	f.checkRecv("d", "e")
}

func newFrame(t *testing.T) frame {
	return frame{
		t:        t,
		send:     NewSend(),
		recv:     NewRecv(),
		received: make([]string, 0),
	}
}

func (f *frame) addAndSend(add ...string) {
	if add == nil {
		add = []string{}
	}
	f.send.AddItems(add...)
	msg, err := f.send.Pack()
	if err != nil {
		if err != ErrNotInited {
			f.t.Error(err)
			return
		}
	}
	items := f.recv.Unpack(msg)
	f.received = append(f.received, items...)
}

func (f *frame) confirm() {
	last := f.recv.LastRecv()
	f.send.Confirm(last)
}

func (f *frame) dropSend() {
	f.send = NewSend()
}

func (f *frame) dropRecv() {
	f.recv = NewRecv()
	f.received = []string{}
}

func (f *frame) checkRecv(wait ...string) {
	if wait == nil {
		wait = []string{}
	}
	if !eqSlices(f.received, wait) {
		f.t.Error("Recv waited ", wait, " got ", f.received)
	}
}
