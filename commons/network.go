package commons

import "encoding/json"

// map[commonPartName]json_PartStruct
type RoomCommonData map[string]string
func (rcd RoomCommonData) Encode() ([]byte, error){
	buf, err:= json.Marshal(rcd)
	if err!=nil{
		Log(LVL_ERROR, "can't encode RoomCommonData")
		return nil, err
	}
	return buf, nil
}

//static method!
func (RoomCommonData) Decode(data []byte) (RoomCommonData, error){
	rcd:=RoomCommonData{}
	err:= json.Unmarshal(data,&rcd)
	if err!=nil{
		Log(LVL_ERROR, "can't decode RoomCommonData")
		return nil, err
	}
	return rcd, nil
}