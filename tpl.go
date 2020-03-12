package snailframe

import (
	"github.com/CloudyKit/jet"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)


type tpl struct {
	tplSet *jet.Set
	tplSuffix string
}

type tplConfig struct {
	Dir string
	Suffix string
	Reload bool
}

func newTpl(cfg tplConfig) *tpl {
	if cfg.Dir == "" {
		cfg.Dir = "template"
	}

	if cfg.Suffix == "" {
		cfg.Suffix = "jet"
	}


	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}

	dir = dir + "/" + cfg.Dir
	tplSuffix := cfg.Suffix

	log.WithFields(log.Fields{
		"Dir":    dir,
		"Suffix": tplSuffix,
		"Debug":cfg.Reload,
	}).Trace("Init Template")

	tplObj := new(tpl)


	tplObj.tplSet = jet.NewHTMLSet(dir)

	isDebug := cfg.Reload
	tplObj.tplSet.SetDevelopmentMode(isDebug)
	tplObj.tplSuffix = tplSuffix
	return tplObj
}

func (this *tpl)AddGlobal(key string, i interface{}) *jet.Set {

	log.WithFields(log.Fields{
		"Global": key,
	}).Trace("Template AddGlobal")

	return this.tplSet.AddGlobal(key, i)
}

/*
简便方法
	var theVar = make(map[string]interface{})
	theVar["xxx"] = "xxx"
	snailrouter.Execute(w,"index.jet",theVar)
*/
func (this *tpl)Execute(w io.Writer, tplName string, maps map[string]interface{}) {
	obj := this.GetTemplate(tplName)

	vars := make(jet.VarMap)

	for mapk, mapv := range maps {
		vars.Set(mapk, mapv)
	}
	obj.Execute(w, vars, nil)
}

func(this *tpl) GetTemplate(Name string) *jet.Template {

	log.WithFields(log.Fields{
		"tplName": Name,
	}).Debug("GetTemplate")

	jt, err := this.tplSet.GetTemplate(Name + "." + this.tplSuffix)
	if err != nil {
		log.WithFields(log.Fields{
			"tplName": Name,
			"err":err,
		}).Panic("GetTemplate Error")
		panic(err)
	}
	return jt
}

func (this *tpl)SetVars(maps map[string]interface{}) jet.VarMap {
	vars := make(jet.VarMap)

	for mapk, mapv := range maps {
		vars.Set(mapk, mapv)
	}

	return vars
}

/*
	templateName := "index.jet"
	t := snailrouter.GetTemplate(templateName)
	theVar := map[string]interface{}{}
	theVar["xxx"] = "xxx"
	vars := snailrouter.SetVars(theVar)
	vars.Set("user", "xx")
	data := r
	t.Execute(w, vars, data.URL.Query());
*/
