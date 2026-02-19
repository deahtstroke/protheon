package transform

import (
	lua "github.com/yuin/gopher-lua"
)

type Transformer interface {
	Transform(input any) (any, error)
}

type LuaTransformer struct {
	Path string
}

func NewLuaTransformer(scriptPath string) *LuaTransformer {
	return &LuaTransformer{
		Path: scriptPath,
	}
}

func (t *LuaTransformer) Transform(input any) (any, error) {
	L := lua.NewState()
	defer L.Close()
	lrecord := goToLua(L, input)

	L.SetGlobal("record", lrecord)

	err := L.DoFile(t.Path)
	if err != nil {
		return nil, err
	}

	output := luaToGo(L.GetGlobal("record")).(map[string]any)
	return output, nil
}

func goToLua(L *lua.LState, data any) lua.LValue {
	switch v := data.(type) {
	case string:
		return lua.LString(v)
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	case []any:
		lt := L.NewTable()
		for _, ele := range v {
			lt.Append(goToLua(L, ele))
		}
		return lt
	case map[string]any:
		lt := L.NewTable()
		for k, v := range v {
			lt.RawSetString(k, goToLua(L, v))
		}
		return lt
	default:
		return lua.LNil
	}
}

func luaToGo(lv lua.LValue) any {
	switch v := lv.(type) {
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case lua.LBool:
		return bool(v)
	case *lua.LTable:
		m := make(map[string]any)
		v.ForEach(func(key, value lua.LValue) {
			m[key.String()] = luaToGo(value)
		})
		return m
	default:
		return nil
	}
}
