package main

import (
	."github.com/Shnifer/magellan/commons"
	"strconv"
	"github.com/Shnifer/magellan/v2"
)

func sphs2sigs(s [15]int) []Signature {
	res := make([]Signature, 0)

	add := func(a, b int) {
		res = append(res, Signature{
			TypeName: strconv.Itoa(a) + "-" + strconv.Itoa(b),
			Dev:      v2.RandomInCircle(1),
		})
	}

	for i, v := range s {
		if v == NONE {
			continue
		}
		switch i {
		case MAGNET, RADIATIONBELT, OXYGEN, OZONE, ION, COREVEL, VULCAN, BIO:
			add(i, v)
		case ATMOMETALS, GASES, PEDOMETALS, COREMADE, LITOMETAL, MIXTURES:
			if v == EARTHANDNEW {
				add(i, EARTH)
				add(i, NEW)
			} else {
				add(i, v)
			}
		case WATER:
			switch v {
			case HARDANDGASOUS:
				add(i, HARD)
				add(i, GASOUS)
			case HARDANDLIQUID:
				add(i, HARD)
				add(i, LIQUID)
			case LIQUIDANDGASOUS:
				add(i, LIQUID)
				add(i, GASOUS)
			case HARDANDLIQUIDANDGASOUS:
				add(i, HARD)
				add(i, LIQUID)
				add(i, GASOUS)
			default:
				add(i, v)
			}
		}
	}

	return res
}

const (
	NONE                   = 0
	WEAK                   = 1
	NORM                   = 2
	STRONG                 = 3
	EARTH                  = 1
	NEW                    = 2
	EARTHANDNEW            = 3
	LIQUID                 = 1
	HARD                   = 2
	GASOUS                 = 3
	HARDANDGASOUS          = 5
	LIQUIDANDGASOUS        = 6
	HARDANDLIQUID          = 7
	HARDANDLIQUIDANDGASOUS = 8
	WAS                    = 4
	RADICAL                = 3
	MOVING                 = 1
	EXTINCT                = 2
	PRESENT                = 1
)

const (
	MAGNET = iota
	RADIATIONBELT
	OXYGEN
	GASES
	ATMOMETALS
	OZONE
	ION
	WATER
	MIXTURES
	PEDOMETALS
	COREMADE
	COREVEL
	VULCAN
	LITOMETAL
	BIO
)

