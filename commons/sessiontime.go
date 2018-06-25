package commons

const sessionTimeElastic = 1

type SessionTime struct {
	target   float64
	elasticT float64
	delta    float64
}

func NewSessionTime(sessionTime float64) *SessionTime {
	return &SessionTime{
		target: sessionTime,
	}
}

func (st *SessionTime) Update(dt float64) {
	st.target += dt
	st.elasticT += dt
}

func (st *SessionTime) Get() (res float64) {
	res = st.target
	if st.elasticT < sessionTimeElastic {
		res += st.delta * (sessionTimeElastic - st.elasticT) / sessionTimeElastic
	}
	return res
}

func (st *SessionTime) MoveTo(sessionTime float64) {
	if sessionTime < st.target {
		return
	}
	if st.target > 0 {
		st.delta = sessionTime - st.Get()
	}
	st.target = sessionTime
}
