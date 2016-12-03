/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose : this file is used as main portal, main count, for sysytem boot application.

* Creation Date : 07-02-2014

* Last Modified : 03-12-2015 11:52:42 PM UTC

* Created By : Kiyor

* Note By: Tim

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"encoding/json"
	"github.com/ccna/ccapi/cclib"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	// 	"github.com/go-xorm/xorm"
	"github.com/kiyor/render"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/sessions"
	htmltemplate "html/template"
	"io/ioutil"
	"log"
	"net/http/pprof"
	//"runtime/pprof""
	"os"
	"syscall"
	texttemplate "text/template"
	// 	"net/http"
	"flag"
	"runtime"
	// 	"sync"
	"os/signal"
	"sync"
	"time"
)

var (
	SUCCESS = "{Code:0,Result:\"Success\"}"
	ac      ApiConfig
	x       cclib.Orm
	m       *martini.Martini
	pushing         = &sync.Mutex{}
	finft   *string = flag.String("p", ":6123", "running port")
	fconf   *string = flag.String("c", "./config.ini", "config file")
)

func init() {
	flag.Parse()
	ac.read(*fconf)
	err := cclib.DBC.Read(*fconf)
	if err != nil {
		panic(err)
	}
	cclib.Logger = cclib.NewLogger(&cclib.LogOptions{
		Name:      "cclib",
		ShowErr:   true,
		ShowDebug: true,
	})
	cclib.InitGit(*fconf)      //init git
	cclib.Verbose = true       //set verbose
	initStoreByFile(*fconf)    //
	os.Setenv("PORT", ac.Port) //
	orm, err := cclib.DBC.InitOrm()
	if err != nil {
		panic(err.Error)
	}
	x = cclib.Orm{DB: orm}

	// martini.Env = martini.Prod
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 	var wg sync.WaitGroup
	// 	wg.Add(1)

	go netserver()
	go run()         //go run time
	go timeChecker() //time checker

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

forever:
	for {
		select {
		case s := <-sig:
			log.Printf("Signal (%d) received, stopping\n", s)
			x.DB.Close()
			log.Println("STOP")
			break forever
		}
	}
	// 	wg.Wait()
}

func netserver() {

}

// cronjob doing something
func timeChecker() { //(dbs []ccnaapi.DB) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			// do nothing
		}
	}
}

// manally run http
func run() { //dbs []ccnaapi.DB) {
	m = martini.New()
	// Setup middleware
	m.Use(logger())
	m.Use(martini.Recovery())
	m.Use(MapEncoder)
	m.Use(gzip.All())

	m.Use(sessions.Sessions("u_session", store))
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
	}))
	// mount static file
	m.Use(martini.Static(ac.Dir+"statics", martini.StaticOptions{
		Prefix:      "statics",
		SkipLogging: true,
	}))

	// setup tmpls folder
	m.Use(render.Renderer(render.Options{
		Directory: ac.Dir + "tmpls",
		Delims:    render.Delims{"[[", "]]"},
		HtmlFuncs: []htmltemplate.FuncMap{
			htmltemplate.FuncMap{"markdown": markDowner},
		},
		TextFuncs: []texttemplate.FuncMap{
			texttemplate.FuncMap{"markdown": markDowner},
		},
	}))
	// Setup end
	// js的 所有gosp 交互 Govlet 中间函数
	r := martini.NewRouter()
	//m := martini.Classic()
	//set debug page
	r.Get("matini/debug/pprof", pprof.Index)
	r.Get("/debug/pprof/cmdline", pprof.Cmdline)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/symbol", pprof.Symbol)
	r.Post("/debug/pprof/symbol", pprof.Symbol)
	r.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	r.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	r.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	// related main.go
	r.Get("/", handlerRootGet)

	//related handler_login.go
	r.Post("/login", handlerLoginPost)

	//related handler_login.go
	r.Get("/login", handlerLoginGet)

	//related handler_logout.go
	r.Get("/logout", handlerLogoutGet)

	//related handler_search.go
	r.Post("/search", handlerSearchPost)

	//related handler_search.go
	r.Get("/search", handlerSearchGet)

	//related handler_configtype.go
	r.Get("/configtype", handlerConfigTypeGet)

	//related main.go
	r.Get("/index", handlerIndexGet)

	//related handler_channels.go
	r.Get("/channels", handlerChannelsGet)

	//related handler_channels.go
	r.Get("/channels.json", handlerChannelsJsonGet)

	//related handler_certs.go
	r.Get("/certs", handlerCertsGet)

	//related handler_editcert.go
	r.Get("/cert/:certid", handlerEditCertGet)

	//related handler_editcert.go
	r.Post("/cert/:certid", handlerEditCertPost)

	//related handler_newcert.go
	r.Get("/newcert", handlerNewCertGet)

	//related handler_newcert.go
	r.Post("/newcert", handlerNewCertPost)

	//related handler_newserver.go
	r.Get("/newserver", handlerNewServerGet)

	//related handler_newserver.go
	r.Post("/newserver", handlerNewServerPost)

	// 	r.Get("/md/:name", handlerMDGet)

	//related handler_md.go
	r.Get("/docs", handlerDocsIndexGet)

	//related handler_md.go
	r.Get("/md/docs/:file", handlerMDDocsGet)

	//related handler_struct.go
	r.Get("/struct/:type/:key", handlerStructGet)

	//related handler_all_type.go
	r.Get("/allacltype", handlerAllAclTypeGet)

	//related handler_all_type.go
	r.Get("/allfeaturetype", handlerAllFeatureTypeGet)

	//handler_servers.go
	r.Get("/servers", handleServersGet)
	//handler_servers.go
	r.Get("/allserver", handlerAllServerGet)
	//handler_servers.go
	r.Get("/server/squid/:server", handlerSquidConfGet)
	//handler_servers.go
	r.Get("/server/origdomain/:server", handlerOrigDomainConfGet)
	//handler_servers.go
	r.Get("/server/domainconf/:server", handlerDomainConfGet)

	//handler_channelconf.go
	r.Get("/channel/:cid", handleChannelConfGet)
	//handler_channelconf.go
	r.Get("/channel/:channel/json", handlerChannelConfJsonGet)
	//handler_channelconf.go
	r.Get("/allgroup", handlerAllGroupGet)
	//handler_channelconf.go
	r.Get("/allcert", handlerAllCertGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/mygroup", handlerChannelGroupGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/link/:gid", handlerChannelLinkToGroupGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/unlink/:gid", handlerChannelUnlinkToGroupGet)
	//handler_channelconf.go
	r.Post("/channel/:cid/basic", handlerChannelBasicPost)
	//handler_channelconf.go
	r.Get("/channel/:cid/commitconf", handlerChannelCommitGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/delete", handlerChannelDeleteGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/log.html", handlerDiffChannelLogGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/histories", handlerChannelHistoriesGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/history/:lid", handlerChannelHistoryGet)
	//handler_channelconf.go
	r.Get("/channel/:cid/rollback/:lid", handlerChannelRollback)
	//handler_channelconf.go
	r.Post("/channel/:cid/historynote/:id", handlerChannelHistoryNoteChangePost)
	//handler_channelconf.go
	r.Post("/channel/:cid/lasthistorynote", handlerChannelLastHistoryNoteChangePost)
	//handler_channelconf.go
	r.Get("/channel/:cid/reload", handlerChannelReload)
	//handler_channelconf.go
	r.Post("/channel/:cid/acl/preview", handlerChannelPreviewAclPost)
	//handler_channelconf.go
	r.Post("/channel/:cid/acl/submit", handlerChannelSubmitAclPost)

	//handler_channelconf.go
	r.Post("/channel/:cid/feature/preview", handlerChannelPreviewFeaturePost)
	//handler_channelconf.go
	r.Post("/channel/:cid/feature/submit", handlerChannelSubmitFeaturePost)

	//handler_channelconf.go
	r.Post("/channel/:cid/rp/preview", handlerChannelPreviewRpPost)
	//handler_channelconf.go
	r.Post("/channel/:cid/rp/submit", handlerChannelSubmitRpPost)

	//handler_channelconf.go
	r.Post("/channel/:cid/nginx/preview", handlerChannelPreviewNginxPost)
	//handler_channelconf.go
	r.Post("/channel/:cid/nginx/submit", handlerChannelSubmitNginxPost)

	//handler_channelconf.go
	r.Post("/channel/:cid/alias/submit", handlerChannelSubmitAlias)
	//handler_channelconf.go
	r.Post("/channel/:cid/route/submit", handlerChannelSubmitRoute)

	//handler_newchannel.go
	r.Get("/channeltemplates", handlerChannelTemplatesGet)
	//handler_newchannel.go
	r.Get("/newchannel", handlerNewChannelGet)
	//handler_newchannel.go
	r.Post("/newchannel", handlerNewChannelPost)

	//handler_channelconf.go
	r.Post("/channel/:cid/domainconf", handlerChannelConfDomainconfPost)

	r.Get("/nodes", handlerNodesGet)
	r.Get("/nodes/deploy", handlerNodeDeployAllGet)
	r.Get("/nodes.json", handlerNodesJsonGet)
	r.Get("/node/:id/deploy", handlerNodeDeployGet)

	/*
		example:
		curl -d {"domain": "changed.com", "origin": "stillorig.com"} http://api/dyn/update/channel/where?id=1
		this func must have admin priv
	*/

	//handler_dynamic_query.go
	r.Post("/dyn/update/:table/where", handlerDynamicUpdateQueryPost)
	/*
		example:
		curl -d changed.com http://api/dyn/update/channel/field/domain/where?id=1
		this func must have admin priv
	*/
	//handler_dynamic_query.go
	r.Post("/dyn/update/:table/field/:field/where", handlerDynamicUpdateSingleQueryPost)

	//handler__testf.go
	r.Post("/test", handlerTestLog)
	//handler__testf.go
	r.Get("/test", handlerTestLog)

	m.Action(r.Handle) //martini载入路由配置
	if *finft != ":6123" {
		m.RunOnAddr(*finft)
	} else {
		m.Run()
	}
}

//get handler root
func handlerRootGet(session sessions.Session, render render.Render) {
	if session.Get(USERNAME) == nil {
		render.Redirect("/login")
		return
	}
	render.Redirect("/index")
}

//
func handlerIndexGet(session sessions.Session, render render.Render) {
	if !isLogin(session) {
		render.Redirect("/login")
		return
	}
	html := initHtml(session)
	md, err := ioutil.ReadFile(ac.Dir + "md/index.md")
	if err != nil {
		log.Println(err.Error())
	}
	html[TITLE] = "CCNA EZCONF"
	html[BODY] = string(md)
	style := []string{"md.css"}
	html[STYLE] = style

	render.TEXT(200, "index", html)
}
