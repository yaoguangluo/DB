/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : handler_channels.go

* Purpose :

* Creation Date : 10-03-2014

* Last Modified : Tue 16 Dec 2014 12:20:37 AM UTC

* Created By : Kiyor

* Note By : Tim
_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"fmt"
	// 	"github.com/ccna/ccapi/cclib"
	"github.com/kiyor/render"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	// 	"strconv"
)

func handlerChannelsGet(r *http.Request, render render.Render, session sessions.Session) {
	if !isLogin(session) {
		render.Redirect("/login")
		return
	}
	// 	res, err := x.GetAllChannel()
	// 	if r.Header.Get("Content-Type") == "application/json" {
	// 		if err != nil {
	// 			render.JSON(200, err.Error())
	// 			return
	// 		}
	// 		render.JSON(200, res)
	// 	}
	html := initHtml(session)
	html[TITLE] = "channel list\n"
	// 	body := "|DOMAIN|ORIGIN|\n"
	// 	body += "|------|------|\n"
	style := []string{"md.css"}
	script := []string{"channels.js"}
	html[STYLE] = style
	html[SCRIPT] = script
	// 	for _, v := range res {
	// 		body += fmt.Sprintf("|[%s](/channel/%d)|[%s](http://%s)|\n", v.Domain, v.Id, v.Origin, v.Origin)
	// 	}
	// 	html[BODY] = body
	// 	log.Println(body)
	render.HTML(200, "channels", html)
	log.Println("!!!", html[USER])
}

func handlerChannelsJsonGet(r *http.Request, render render.Render, session sessions.Session) {
	if !isLogin(session) {
		return
	}
	res, err := x.GetAllChannel()
	if err != nil {
		render.JSON(200, err.Error())
		return
	}
	render.JSON(200, res)
}
