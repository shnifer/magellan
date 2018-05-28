package storage

import (
	"encoding/json"
	"log"
)

type Request struct {
	INeedFullKey string `json:"n"`
}

type Responce struct {
	IHaveFullKeys []string `json:"h"`
	YourKeyVal    string   `json:"v"`
}

func (r Request) Encode() []byte {
	buf, err := json.Marshal(r)
	if err != nil {
		log.Panicln("Can't marshal request", err)
	}
	return buf
}

func (Request) Decode(buf []byte) (r Request, err error) {
	err = json.Unmarshal(buf, &r)
	if err != nil {
		log.Println("Can't unmarshal request", err)
		return Request{}, err
	}
	return r, nil
}

func (r Responce) Encode() []byte {
	buf, err := json.Marshal(r)
	if err != nil {
		log.Panicln("Can't marshal request", err)
	}
	return buf
}

func (Responce) Decode(buf []byte) (r Responce, err error) {
	err = json.Unmarshal(buf, &r)
	if err != nil {
		log.Println("Can't unmarshal request", err)
		return Responce{}, err
	}
	return r, nil
}
