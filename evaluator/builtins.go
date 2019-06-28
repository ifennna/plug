package evaluator

import "plug/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid number of arguments to `len`, expected 1, got %d", len(args))
			}

			switch argument := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(argument.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
