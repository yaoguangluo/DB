/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : logger.go

* Purpose :this file is talke about the system logger show.

* Creation Date : 08-19-2014

* Last Modified : 03-13-2015 06:31:15 PM UTC

* Created By : Kiyor

* note By : Tim

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"github.com/go-martini/martini"
	// 	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
	"time"
)

// func logger(session sessions.Session) martini.Handler {
func logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}
		// 		var user string
		// 		user = session.Get("username").(string)

		// 		log.Printf("Started %s %s for %s by %s", req.Method, req.URL.Path, addr, user)
		log.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(martini.ResponseWriter) // cal martini response

		c.Next()

		// 		log.Printf("Completed %v %s in %v by %s", rw.Status(), http.StatusText(rw.Status()), time.Since(start), user)
		log.Printf("Completed %s %s for %s [%v] %s in %v", req.Method, req.URL.Path, addr, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
