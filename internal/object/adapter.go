package object

import (
	"fmt"

	"github.com/botscubes/bot-components/context"
)

func ConvertContextToEnv(ctx *context.Context, env *Env, vars *[]string) (*Env, error) {
	for _, varName := range *vars {
		value, ok := ctx.GetRawValue(varName)
		if !ok {
			return nil, fmt.Errorf("variable does not exists")
		}

		switch v := value.(type) {
		case int, int64:
			env.Set(varName, &Integer{Value: v.(int64)})
		case float32, float64:
			env.Set(varName, &Integer{Value: int64(v.(float64))})
		case string:
			env.Set(varName, &String{Value: v})
		case bool:
			env.Set(varName, &Boolean{Value: v})
		case map[string]any:
			hm, err := convertMapToHashMap(v)
			if err != nil {
				return nil, fmt.Errorf("hashmap convert error")
			}

			env.Set(varName, hm)
		case []any:
			arrayObj, err := convertArray(v)
			if err != nil {
				return nil, err
			}
			env.Set(varName, arrayObj)
		default:
			return nil, fmt.Errorf("неизвестный тип данных: %T", v)
		}
	}
	return env, nil
}

func convertMapToHashMap(data map[string]any) (Object, error) {
	pairs := make(map[HashKey]HashPair)

	for k, v := range data {
		var keyObject Object
		var valueObject Object

		keyObject = &String{Value: k}
		hashKey, ok := keyObject.(Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %T", hashKey)
		}

		switch val := v.(type) {
		case int, int64:
			valueObject = &Integer{Value: val.(int64)}
		case float32, float64:
			valueObject = &Integer{Value: int64(val.(float64))}
		case string:
			valueObject = &String{Value: val}
		case bool:
			valueObject = &Boolean{Value: val}
		case map[string]any:
			innerHashMap, err := convertMapToHashMap(val)
			if err != nil {
				return nil, err
			}
			valueObject = innerHashMap
		default:
			return nil, fmt.Errorf("неизвестный тип данных в hashmap: %T", val)
		}

		pairs[hashKey.HashKey()] = HashPair{Key: keyObject, Value: valueObject}
	}

	return &HashMap{Pairs: pairs}, nil
}

func convertArray(data []any) (Object, error) {
	elements := make([]Object, len(data))

	for i, v := range data {
		var elementObject Object
		switch val := v.(type) {
		case int, int64:
			elementObject = &Integer{Value: val.(int64)}
		case float32, float64:
			elementObject = &Integer{Value: int64(val.(float64))}
		case string:
			elementObject = &String{Value: val}
		case bool:
			elementObject = &Boolean{Value: val}
		case map[string]any:
			hashMapObj, err := convertMapToHashMap(val)
			if err != nil {
				return nil, err
			}
			elementObject = hashMapObj
		case []any:
			arrayObj, err := convertArray(val)
			if err != nil {
				return nil, err
			}
			elementObject = arrayObj
		default:
			return nil, fmt.Errorf("неизвестный тип данных в массиве: %T", val)
		}
		elements[i] = elementObject
	}

	return &Array{Elements: elements}, nil
}

func ExtractRawValueFromObject(obj Object) (any, bool) {
	switch obj.Type() {
	case INTEGER_OBJ:
		return obj.(*Integer).Value, true
	case BOOLEAN_OBJ:
		return obj.(*Boolean).Value, true
	case STRING_OBJ:
		return obj.(*String).Value, true
	case ARRAY_OBJ:
		a := make([]any, len(obj.(*Array).Elements))
		for id, e := range obj.(*Array).Elements {
			v, ok := ExtractRawValueFromObject(e)
			if !ok {
				return v, false
			}

			a[id] = v
		}
		return a, true
	case HASH_MAP_OBJ:
		m := make(map[string]any)
		for _, el := range obj.(*HashMap).Pairs {
			v, ok := ExtractRawValueFromObject(el.Value)
			if !ok {
				return v, false
			}

			m[el.Key.ToString()] = v
		}

		return m, true
	case NULL_OBJ:
		return nil, true
	case ERROR_OBJ:
		return obj.ToString(), false
	default:
		return fmt.Sprintf("unsupported result type: %s", obj.Type()), false
	}
}
