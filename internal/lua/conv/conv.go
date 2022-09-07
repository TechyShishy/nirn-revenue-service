package conv

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func String(l lua.LValue) (string, error) {
	lString, ok := l.(lua.LString)
	if !ok {
		return "", fmt.Errorf("wanted string, got %v from lua: %#v", l.Type(), l)
	}
	return string(lString), nil
}

func Int(l lua.LValue) (int, error) {
	lNumber, ok := l.(lua.LNumber)
	if !ok {
		return 0, fmt.Errorf("wanted number, got %v from lua: %#v", l.Type(), l)
	}
	lFloat := float64(lNumber)
	lInt := int(lFloat)
	if lFloat != float64(lInt) {
		return 0, fmt.Errorf("wanted int, got float from lua: %#v", lFloat)
	}
	return lInt, nil
}

func Uint(l lua.LValue) (uint, error) {
	lInt, err := Int(l)
	if err != nil {
		return 0, err
	}
	lUint := uint(lInt)
	if lInt != int(lUint) {
		return 0, fmt.Errorf("wanted uint, got int from lua: %#v", lInt)
	}
	return lUint, nil
}

func Bool(l lua.LValue) (bool, error) {
	lBool, ok := l.(lua.LBool)
	if !ok {
		return false, fmt.Errorf("wanted bool, got %v from lua: %#v", l.Type(), l)
	}
	return bool(lBool), nil
}

func Table(l lua.LValue) (*lua.LTable, error) {
	lTable, ok := l.(*lua.LTable)
	if !ok {
		return &lua.LTable{}, fmt.Errorf("wanted lua.LTable, got %v from lua: %#v", l.Type(), l)
	}
	return lTable, nil
}
