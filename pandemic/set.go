package pandemic

type Set map[string]struct{}

func Init(ks []string) Set {
	s := Set{}
	for _, k := range ks {
		s[k] = struct{}{}
	}
	return s
}

func (s Set) Contains(k string) bool {
	_, ok := s[k]
	return ok
}

func (s Set) Add(k string) Set {
	s[k] = struct{}{}
	return s
}

func (s Set) Remove(k string) (Set, bool) {
	if _, ok := s[k]; !ok {
		return s, false
	}
	delete(s, k)
	return s, true
}

func (s Set) Size() int {
	return len(s)
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
