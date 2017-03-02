package common

func ContainString(sl []string, v string) (int, bool) {
	for index, vv := range sl {
		if vv == v {
			return index, true
		}
	}
	return 0, false
}
