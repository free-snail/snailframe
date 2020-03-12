package snailframe

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)
type FrameWork struct {
	DataCon *dbConn
	Tpl *tpl
	FrameRouter *SnailRouter
	RouterMaps RouterMap
	Config coreConfig
}

type coreConfig struct {
	ListenAddr string
	Mysql dbConfig
	Log logConfig
	Tpl tplConfig
	Res rdataConfig
}

func NewFrame() *FrameWork {

	//初始化配置
	var config coreConfig
	NewConf(&config,"config.toml")

	//初始化日志
	initLog(config.Log)

	myFrameWork := &FrameWork{}
	//初始化数据库
	myFrameWork.DataCon = newDB(config.Mysql)

	//初始化模板
	myFrameWork.Tpl = newTpl(config.Tpl) //为什么还要开放出去？是为了能够执行addGlobal等增加全局函数等操作

	//初始化路由 将模板传进去
	myFrameWork.FrameRouter = newRouter(config.Res,myFrameWork.Tpl,myFrameWork.DataCon)

	myFrameWork.Config = config

	return myFrameWork
}



/**
var routerMap = map[string]router.RegFunc{
	"/":                  ctl.Index,
	"/poet/":             ctl.Poet,
	"/poet/info/{aid}":   ctl.PoetInfo,
	"/poetry/":           ctl.Poetry,
	"/poetry/info/{sid}": ctl.PoetryInfo,
	"/whatnews/":         ctl.Whatnews,
}
 */
func (this *FrameWork)StartServer(){

	if this.RouterMaps == nil {
		this.FrameRouter.HandleFunc("/",Welcome)
	}else{
		for k,v := range this.RouterMaps {
			log.WithFields(log.Fields{"Path": k,"Func": v,}).Trace("Load RouterMap")
			this.FrameRouter.HandleFunc(k,v)
		}
	}

	ListenAddr := this.Config.ListenAddr
	log.WithFields(log.Fields{"addr":  ListenAddr,}).Info("Serve Start Listen...")

	/*server := &http.Server{Addr: ListenAddr, Handler: this.FrameRouter,ReadTimeout:time.Duration(time.Second*1),
		ReadHeaderTimeout:time.Duration(time.Second*1),WriteTimeout:time.Duration(time.Second*1),
		IdleTimeout:time.Duration(time.Second*1),
	}
	server.ListenAndServe()*/
	//return server.ListenAndServe()

	err := http.ListenAndServe(ListenAddr, this.FrameRouter) //设置监听的端口
	if err != nil {
		log.WithFields(log.Fields{"err":  err,}).Fatal("Cant Listen...")
	}

}






