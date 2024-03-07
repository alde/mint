package object

import "hash/fnv"

type String struct {
	Value string
}

func (b *String) Inspect() string  { return b.Value }
func (b *String) Type() ObjectType { return STRING_OBJ }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
