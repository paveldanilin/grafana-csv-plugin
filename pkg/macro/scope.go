package macro

type Scope struct {
	Variables map[string]interface{}
}

func NewScope() *Scope {
	return &Scope{
		Variables: make(map[string]interface{}),
	}
}

func (s *Scope) HasVar(name string) bool {
	_, has := s.Variables[name]
	return has
}

func (s *Scope) GetVar(name string) interface{} {
	if s.HasVar(name) {
		return s.Variables[name]
	}
	return nil
}

func (s *Scope) SetVar(name string, value interface{}) {
	s.Variables[name] = value
}
