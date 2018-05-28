package storage

import (
	"errors"
	"strings"
)

const separator = "~"

type ObjectKey struct {
	Glyph, Area, Node, Key string
}

// strings must not have "/","\" or "~" symbols
func newKey(glyph, area, node, key string) ObjectKey {
	return ObjectKey{
		Glyph: glyph,
		Area:  area,
		Node:  node,
		Key:   key,
	}
}

func ReadKey(fullKey string) (k ObjectKey, err error) {
	parts := strings.Split(fullKey, separator)
	if len(parts) != 3 {
		return ObjectKey{}, errors.New("splitKey: len(parts)!=3")
	}
	k.Area, k.Node, k.Key = parts[0], parts[1], parts[2]
	if k.Area == "" || k.Node == "" || k.Key == "" {
		return ObjectKey{}, errors.New("splitKey: some fields are empty " + fullKey)
	}
	switch k.Area[:1] {
	case glyphDel:
		k.Glyph = glyphDel
		k.Area = k.Area[1:]
	}
	return k, nil
}

func (k ObjectKey) fullKey() string {
	return k.Glyph + strings.Join([]string{k.Area, k.Node, k.Key}, separator)
}
