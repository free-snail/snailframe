package snailframe

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)


type RouterMap map[string]func(*RData)

type SnailRouter struct {
	*mux.Router
	config rdataConfig
	tpl *tpl
	dbconn *dbConn
}

func newRouter(cfg rdataConfig,Tpl *tpl,db *dbConn) *SnailRouter {
	MuxRouter := mux.NewRouter().SkipClean(false).UseEncodedPath().StrictSlash(true)

	//https://github.com/gorilla/mux/blob/master/mux.go

	myRouter := &SnailRouter{}
	myRouter.Router = MuxRouter

	myRouter.Router.NotFoundHandler = http.HandlerFunc(myRouter.notFund)
	MuxRouter.Use(myRouter.LogInfo())

	myRouter.config = cfg
	myRouter.tpl = Tpl
	myRouter.dbconn = db

	log.WithFields(log.Fields{
		"SafetheUrl":  cfg.SafetheUrl,
		"SessionKey": cfg.SessionKey,
		"SessionMaxAge": cfg.SessionMaxAge,
		"CookiePath": cfg.CookiePath,
		"CookieDomain": cfg.CookieDomain,
		"CookieSameSite":  cfg.CookieSameSite,
		"CookieSecure": cfg.CookieSecure,
		"CookieHttpOnly": cfg.CookieHttpOnly,
	}).Trace("Init router...")

	return myRouter
}

func (this *SnailRouter)HandleFunc(path string, f func(*RData)) *mux.Route {

	//这里采用闭包，会不会有问题？
	reFunc := func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			//显示错误信息
			if this.config.ShowErro {
				if pr := recover(); pr != nil {
					this.showError(w,pr, debug.Stack())
				}
			}
		}()

		rdata := new(RData)
		rdata.ResponseWriter = w
		rdata.Request = r
		rdata.config = this.config
		rdata.tpl = this.tpl
		rdata.dbconn = this.dbconn

		f(rdata)
	}

	return this.Router.HandleFunc(path, reFunc)
}



func (this *SnailRouter)showError(w http.ResponseWriter,recoverInfo interface{},bugtrace []byte){
	info := []byte(fmt.Sprintf("<html>panic recover=====: %v\r\n<br /><br /><hr /><pre>", recoverInfo))
	bugInfo := append(info, debug.Stack()...)
	bugInfo = append(bugInfo, []byte("</pre></html>")...)
	w.WriteHeader(501)
	w.Write(bugInfo)
}

func (this *SnailRouter)notFund(w http.ResponseWriter, r *http.Request){
	log.WithFields(log.Fields{
		"Url":  r.URL.String(),
	}).Info("404 Not found")
	w.WriteHeader(404)
	w.Write([]byte("404 not found"))
}


func (this *SnailRouter)LogInfo() mux.MiddlewareFunc {

	return func(f http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log.WithFields(log.Fields{
				"Url":  r.URL.String(),
			}).Info("An request is come")

			f.ServeHTTP(w, r)
		})
	}
}

