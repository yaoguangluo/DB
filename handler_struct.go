/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : handler_struct.go

* Purpose :

* Creation Date : 10-02-2014

* Last Modified : Thu 04 Dec 2014 09:30:40 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"github.com/ccna/ccapi/cclib"
	"github.com/go-martini/martini"
	"github.com/kiyor/render"
	"github.com/martini-contrib/sessions"
)

func handlerStructGet(session sessions.Session, params martini.Params, render render.Render) {
	typ := params["type"]
	key := params["key"]
	switch typ {
	case "acl":
		render.JSON(200, cclib.NewAcl(key))
		return
	case "feature":
		render.JSON(200, cclib.NewFeature(key))
		return
	case "rp":
		render.JSON(200, cclib.NewRefreshPattern())
		return
	case "channel":
		switch key {
		case "alias":
			render.JSON(200, cclib.NewChannelAlias())
			return
		case "route":
			render.JSON(200, cclib.NewCustomRoute())
			return
		}
	}
}
