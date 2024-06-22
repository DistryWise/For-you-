package scope

type Scope = map[string]string

func NewScope() Scope {
	scopeEl := make(Scope)
	return scopeEl
}
