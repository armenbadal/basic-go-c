package parser

type scope struct {
	parent *scope
	items  []string
}

type symbols struct {
	scopes *scope
}

func (s *symbols) openScope() {
	s.scopes = &scope{s.scopes, make([]string, 0, 8)}
}

func (s *symbols) closeScope() {
	s.scopes = s.scopes.parent
}

func (s *symbols) add(n string) {
	s.scopes.items = append(s.scopes.items, n)
}

func (s *symbols) find(n string) bool {
	p := s.scopes
	for p != nil {
		for _, e := range p.items {
			if e == n {
				return true
			}
		}
		p = p.parent
	}

	return false
}
