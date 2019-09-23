package evaluator

import (
	"github.com/noculture/plug/object"
)

var Out func(value string)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{Function: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("invalid number of arguments to `len`, expected 1, got %d", len(args))
		}

		switch argument := args[0].(type) {
		case *object.Array:
			return &object.Integer{Value: int64(len(argument.Elements))}
		case *object.String:
			return &object.Integer{Value: int64(len(argument.Value))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	},
	},
	"first": &object.Builtin{Function: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("invalid number of arguments, expected 1, got %d", len(args))
		}

		if args[0].Type() != object.ARRAY {
			return newError("argument to `first` must be an array, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}

		return NULL
	}},
	"last": &object.Builtin{Function: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("invalid number of arguments, expected 1, got %d", len(args))
		}

		if args[0].Type() != object.ARRAY {
			return newError("argument to `last` must be an array, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		if length > 0 {
			return arr.Elements[length-1]
		}

		return NULL
	}},
	"rest": &object.Builtin{Function: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("invalid number of arguments, expected 1, got %d", len(args))
		}

		if args[0].Type() != object.ARRAY {
			return newError("argument to `rest` must be an array, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		if length > 0 {
			newElements := make([]object.Object, length-1, length-1)
			copy(newElements, arr.Elements[1:length])
			return &object.Array{Elements: newElements}
		}

		return NULL
	}},
	"push": &object.Builtin{Function: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("invalid number of arguments, expected 2, got %d", len(args))
		}

		if args[0].Type() != object.ARRAY {
			return newError("first argument to `push` not supported, expected ARRAY, got %s", args[0].Type())
		}

		arr := args[0].(*object.Array)
		length := len(arr.Elements)

		newElements := make([]object.Object, length+1, length+1)
		copy(newElements, arr.Elements)
		newElements[length] = args[1]

		return &object.Array{Elements: newElements}
	}},
	"print": &object.Builtin{Function: func(args ...object.Object) object.Object {
		for _, arg := range args {
			Out(arg.Inspect())
		}
		return NULL
	},
	},
}
