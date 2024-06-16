package evaluator

import (
	"strconv"

	"github.com/botscubes/bql/internal/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments: %d want: 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("type of argument not supported: %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments: %d want: 2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument must be ARRAY, got: %s", args[0].Type())
			}

			args[0].(*object.Array).Elements = append(args[0].(*object.Array).Elements, args[1])
			return args[0]
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments: %d want: 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument must be ARRAY, got: %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments: %d want: 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument must be ARRAY, got: %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}

			return NULL
		},
	},
	"intToString": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments: %d want: 1", len(args))
			}

			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument must be INTEGER, got: %s", args[0].Type())
			}

			number := args[0].(*object.Integer).Value

			return &object.String{Value: strconv.FormatInt(number, 10)}
		},
	},
}
