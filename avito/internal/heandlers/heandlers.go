package heandler

import "github.com/julienschmidt/httprouter"

type Heandlers interface {
	Register(router *httprouter.Router)
}
