package main

type mediCounters struct {
	wasLow          bool
	wasHigh         bool
	counterMidTotal int
	counterMid      int
	counterLow      int
	levels          [3]float64
	bioInfSignals   [3][7]int
	nucleoSignals   [3][7]int
}

func newCounter(o CounterOpts) *mediCounters {
	return &mediCounters{
		levels:        o.Levels,
		bioInfSignals: o.BioInf,
		nucleoSignals: o.Nucleo,
	}
}

func (c *mediCounters) AddValue(x float64) {
	//count medium timers if needed before early returns
	if c.counterMidTotal > 0 {
		c.counterMidTotal++
		if c.counterMidTotal > DEFVAL.MediMidTotalS {
			if c.counterMid > DEFVAL.MediMidNeededS {
				c.sendSignal(1)
			}
			c.counterMid = 0
			c.counterMidTotal = 0
		}
	}

	//no effect
	if x < c.levels[0] {
		return
	}

	//low effect
	if x < c.levels[1] {
		if !c.wasLow {
			c.counterLow++
			if c.counterLow > DEFVAL.MediLowCounterS {
				c.wasLow = true
				c.sendSignal(0)
			}
		}
		return
	}

	//mid signal
	if x < c.levels[2] {
		if c.counterMidTotal == 0 {
			c.counterMidTotal++
		}
		c.counterMid++
		return
	}

	//hard effect
	if !c.wasHigh {
		c.wasHigh = true
		c.sendSignal(2)
	}
}

func (c *mediCounters) sendSignal(lvl int) {
	go sendAlice(c.bioInfSignals[lvl], c.nucleoSignals[lvl])
}
