package evalutator

import (
	"fmt"

	"alde.nu/mint/object"
)

func lengthFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments to `len`. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}

func typeFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments to `type`. got=%d, want=1", len(args))
	}

	return &object.String{Value: string(args[0].Type())}
}

func firstFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)
	if len(arr.Elements) > 0 {
		return arr.Elements[0]
	}
	return NULL
}

func lastFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		return arr.Elements[length-1]
	}
	return NULL
}

func restFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}
	return NULL
}

func pushFn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}
	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]
	return &object.Array{Elements: newElements}
}

func putsFn(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return NULL
}

var builtins = map[string]*object.Builtin{
	"len":   {Fn: lengthFn},
	"type":  {Fn: typeFn},
	"first": {Fn: firstFn},
	"last":  {Fn: lastFn},
	"rest":  {Fn: restFn},
	"push":  {Fn: pushFn},
	"puts":  {Fn: putsFn},
}
