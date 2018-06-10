package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/static"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"io/ioutil"
	"log"
)

type Signature struct {
	TypeName  string
	SigString string
	//deviation of this instance
	//supposed to be Len<=1
	Dev v2.V2
}

const sigAtlasFN = "signatures.json"
const sigParticlesFN = "particles.json"

const (
	SIG_SPAWNPERIOD = "SpawnPeriod"
	SIG_LIFETIME    = "LifeTime"
	SIG_VELSPAWN    = "VelAndSpawnStr"
	SIG_ANGF        = "AngStr"
	SIG_SIZEF       = "SizeStr"
	SIG_ALPHAF      = "AlphaStr"
)

var signatures SignatureAtlas
var particles SignatureParticleAtlas

type SignatureAtlas map[string]SignatureType
type SignatureParticleAtlas map[string]SignatureParticle

type SignatureParticle struct {
	SpriteName   string
	DoRandomLine bool
	FPS          float64
	CycleType    int
}

type SignatureType struct {
	ParticleName   string
	Color          color.RGBA
	SpawnPeriod    float64
	LifeTime       float64
	VelAndSpawnStr string
	AngStr         string
	SizeStr        string
	AlphaStr       string
	ApplyDevOn     []string
	applyDevOn     map[string]bool
}

func (s Signature) Type() SignatureType {
	return signatures[s.TypeName]
}
func (s Signature) Particle() SignatureParticle {
	return s.Type().Particle()
}
func (st SignatureType) Particle() SignatureParticle {
	return particles[st.ParticleName]
}

func (s Signature) DevV(name string) v2.V2 {
	if s.Type().applyDevOn[name] {
		return s.Dev
	} else {
		return v2.ZV
	}
}

func (s Signature) DevF(name string) float64 {
	return s.DevV(name).LenSqr()
}

func (s Signature) DevK(name string, widePercent float64) float64 {
	k := s.DevF(name)*2 - 1
	return 1 + (k * widePercent / 100)
}

func InitSignatureAtlas() {
	saveSignatureExample("example_" + sigAtlasFN)
	saveParticleExample("example_" + sigParticlesFN)
	var data []byte
	var err error
	data, err = static.Load("signatures", sigParticlesFN)
	if err != nil {
		panic("Can't find particle atlas file " + sigParticlesFN)
	}
	particles = make(SignatureParticleAtlas)
	err = json.Unmarshal(data, &particles)
	if err != nil {
		panic(err)
	}
	data, err = static.Load("signatures", sigAtlasFN)
	if err != nil {
		panic("Can't find signature atlas file " + sigAtlasFN)
	}
	signatures = make(SignatureAtlas)
	err = json.Unmarshal(data, &signatures)
	if err != nil {
		panic(err)
	}
	for name, v := range signatures {
		if _, ok := particles[v.ParticleName]; !ok {
			log.Panicln("Signature", name, "particle", v.ParticleName, "not found")
		}
		v.applyDevOn = make(map[string]bool, len(v.ApplyDevOn))
		for _, str := range v.ApplyDevOn {
			v.applyDevOn[str] = true
		}
		signatures[name] = v
	}
}

func saveSignatureExample(fn string) {
	exAtlas := make(SignatureAtlas)
	exAtlas["name"] = SignatureType{
		ParticleName:   "particleName",
		SpawnPeriod:    1,
		LifeTime:       2,
		VelAndSpawnStr: "funcName",
		AngStr:         "funcName",
		SizeStr:        "funcName",
		AlphaStr:       "funcName",
		ApplyDevOn:     []string{SIG_SPAWNPERIOD, SIG_LIFETIME, SIG_VELSPAWN, SIG_ANGF, SIG_SIZEF, SIG_ALPHAF},
	}

	buf, err := json.Marshal(exAtlas)
	if err != nil {
		panic(err)
	}
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, buf, "", "  ")
	err = ioutil.WriteFile(fn, identbuf.Bytes(), 0)
	if err != nil {
		panic("can't write texture atlas example " + err.Error())
	}
}

func saveParticleExample(fn string) {
	exAtlas := make(SignatureParticleAtlas)
	exAtlas["name"] = SignatureParticle{
		SpriteName:   "texAtlasName",
		DoRandomLine: true,
		FPS:          20,
		CycleType:    1,
	}

	buf, err := json.Marshal(exAtlas)
	if err != nil {
		panic(err)
	}
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, buf, "", "  ")
	err = ioutil.WriteFile(fn, identbuf.Bytes(), 0)
	if err != nil {
		panic("can't write texture atlas example " + err.Error())
	}
}
