package hoq

type Headers map[string][]string

func (h Headers) Serialize() string {
	out := ""
	//todo check omit values after the first one
	for name, values := range h {
		if len(values) == 0 {
			continue
		}
		out += name + ": " + values[0]
	}
	return out
}
