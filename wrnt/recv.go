package wrnt

type Recv struct {
	lastN int
}

func NewRecv() *Recv {
	return &Recv{
		lastN: 0,
	}
}

func (r *Recv) Unpack(msg Storage) []string {
	var firstInd int
	if r.lastN >= msg.BaseN {
		firstInd = r.lastN - msg.BaseN + 1
	}
	if firstInd < 0 {
		firstInd = 0
	}
	r.lastN = msg.BaseN + len(msg.Items) - 1
	return msg.Items[firstInd:]
}

func (r *Recv) LastRecv() int {
	return r.lastN
}
