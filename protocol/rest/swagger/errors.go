package swagger

import "github.com/golang-tire/pkg/types"

var (
	errDef    types.MapSI
	errFldDef types.MapSI
	err400    = types.MapSI{
		"description": "Input error",
		"schema": types.MapSI{
			"$ref": "#/definitions/ErrorFldResponse",
		},
	}
	err401 = types.MapSI{
		"description": "Returned when not authenticated",
		"schema": types.MapSI{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
	err403 = types.MapSI{
		"description": "Returned when not authorized",
		"schema": types.MapSI{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
	err404 = types.MapSI{
		"description": "Returned when the route is not correct",
		"schema": types.MapSI{
			"$ref": "#/definitions/ErrorResponse",
		},
	}
)
