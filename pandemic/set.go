package pandemic

type Set map[string]bool

func Init(ks []string) Set {
	s := Set{}
	for _, k := range ks {
		s[k] = true
	}
	return s
}

func (s Set) Contains(k string) bool {
	return s[k] == true
}

func (s Set) Add(k string) Set {
	s[k] = true
	return s
}

func Intersection(s1 Set, s2 Set) Set {
	s3 := Set{}

	for k, _ := range s1 {
		if s2.Contains(k) {
			s3.Add(k)
		}
	}

	return s3
}
