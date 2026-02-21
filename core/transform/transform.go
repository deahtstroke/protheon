package transform

import (
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

type Transformer interface {
	Transform(input any) (any, error)
	Close() error
}

type LuaTransformer struct {
	LuaVM *lua.LState
	Path  string
}

func NewLuaTransformer(scriptPath string) *LuaTransformer {
	L := lua.NewState()
	return &LuaTransformer{
		Path:  scriptPath,
		LuaVM: L,
	}
}

func (t *LuaTransformer) Transform(input any) (any, error) {
	lrecord := goToLua(t.LuaVM, input)

	t.LuaVM.SetGlobal("record", lrecord)

	err := t.LuaVM.DoFile(t.Path)
	if err != nil {
		return nil, err
	}

	output := luaToGo(t.LuaVM.GetGlobal("record")).(map[string]any)
	return output, nil
}

func (t *LuaTransformer) Close() error {
	t.LuaVM.Close()
	return nil
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
		rv := reflect.ValueOf(data)

		if rv.Kind() == reflect.Pointer {
			if rv.IsNil() {
				return lua.LNil
			}
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Struct:
			lt := L.NewTable()
			rt := rv.Type()

			for i := 0; i < rv.NumField(); i++ {
				field := rt.Field(i)

				if field.PkgPath != "" {
					continue
				}

				lt.RawSetString(field.Name, goToLua(L, rv.Field(i).Interface()))
			}
			return lt
		case reflect.Slice, reflect.Array:
			lt := L.NewTable()
			for i := range rv.Len() {
				lt.Append(goToLua(L, rv.Index(i).Interface()))
			}
			return lt
		}

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
