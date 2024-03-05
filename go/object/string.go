package object

type String struct {
	Value string
}

func (b *String) Inspect() string  { return b.Value }
func (b *String) Type() ObjectType { return STRING_OBJ }
