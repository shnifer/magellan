package input

const historyLen = 6
var axisBuffers map[string][]float64

func init(){
	axisBuffers = make(map[string][]float64)
}

func bufferAndGet(name string, v float64) float64{
	b := append(axisBuffers[name], v)
	cut:=len(b)-historyLen
	if cut>0{
		b = b[cut:]
	}
	axisBuffers[name] = b
	var s float64
	for _,v:=range b{
		s+=v
	}
	return s/float64(len(b))
}