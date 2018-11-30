package main

import (
	"encoding/json"
	"fmt"
	"github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/log"
	"net/http"
	"strings"
)

func loadHyBoostList() map[string]BoostParams {
	res := make(map[string]BoostParams)
	buf, err := commons.DoReq(http.MethodGet, DEFVAL.RequestHyBoostListAddr, []byte{})
	if err != nil {
		Log(LVL_ERROR, "can't load HyBoostList, addr ", DEFVAL.RequestHyBoostListAddr, " err: ", err)
		return res
	}
	var list []BoostParams
	err = json.Unmarshal(buf, &list)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal HyBoostList, err ", err)
		return res
	}
	for _, v := range list {
		res[normBoostN(v.Password)] = v
	}
	return res
}

type ReqBoostResp struct {
	Status string `json:"status"`
}

func ReqHyBoost(pass string) bool {
	reqBody := fmt.Sprintf(`{"password": "%v"}`, pass)
	buf, err := commons.DoReq(http.MethodPost, DEFVAL.RequestHyBoostUseAddr, []byte(reqBody))
	if err != nil {
		LogGame("failedReqs", false, "can't do ReqHyBoost, err: ", err)
		return false
	}
	var resp ReqBoostResp
	err = json.Unmarshal(buf, resp)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal ReqHyBoost, err: ", err)
		return false
	}

	return resp.Status == "ok"
}

func (s *engiScene) tryBoost(boostPass string) bool {
	boostK := normBoostN(boostPass)
	bp, ok := boostList[boostK]
	if !ok {
		return false
	}
	var sysN int
	switch bp.NodeType {
	case "march_engine":
		sysN = 0
	case "shunter":
		sysN = 1
	case "warp_engine":
		sysN = 2
	case "shields":
		sysN = 3
	case "radar":
		sysN = 4
	case "scaner":
		sysN = 5
	case "fuel_tank":
		sysN = 6
	case "lss":
		sysN = 7
	default:
		Log(LVL_ERROR, "Unknown Hy boost nodeType: ", bp.NodeType)
		return false
	}

	lifeTime := bp.BaseTime + bp.AZBonus*Data.EngiData.AZ[sysN]/100
	boost := commons.Boost{
		LeftTime: lifeTime,
		Power:    bp.BoostPower,
		SysN:     sysN,
	}

	used := ReqHyBoost(bp.Password)
	if !used {
		return false
	}

	delete(boostList, boostK)
	Data.EngiData.Boosts = append(Data.EngiData.Boosts, boost)
	s.doTargetAZDamage(sysN, bp.AZDmg)
	return true
}

func updateBoosts(dt float64) {
	for i, v := range Data.EngiData.Boosts {
		lt := v.LeftTime
		if lt > 0 {
			lt -= dt
			if lt < 0 {
				lt = 0
			}
			Data.EngiData.Boosts[i].LeftTime = lt
		}
	}
}

func normBoostN(s string) string {
	res := ""
	digits := "1234567890"
	for _, c := range s {
		if strings.ContainsRune(digits, c) {
			res += string(c)
		}
	}
	return res
}
