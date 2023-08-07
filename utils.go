package flow

import "strconv"

func AsBool(v any) bool {
	if v == nil {
		return false
	}

	switch vt := v.(type) {
	default:
		return false
	case bool:
		return vt
	case string:
		// True string values
		// 1, t, T, TRUE, true, True
		// False string values
		// 0, f, F, FALSE, false, False
		vb, err := strconv.ParseBool(vt)
		if err != nil {
			panic(err)
		}
		return vb
	case int:
		return vt != 0
	case *bool:
		if vt == nil {
			return false
		}
		return *vt
	case *string:
		if vt == nil {
			return false
		}
		// True string values
		// 1, t, T, TRUE, true, True
		// False string values
		// 0, f, F, FALSE, false, False
		vb, err := strconv.ParseBool(*vt)
		if err != nil {
			panic(err)
		}
		return vb
	case *int:
		if vt == nil {
			return false
		}
		return *vt != 0
	}
}
