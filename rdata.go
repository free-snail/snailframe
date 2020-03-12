// 包括ReQuestData 与 RePonseData的处理
package snailframe

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"snailframe/utils/convert"
	"snailframe/utils/mime"
	"time"
)

type rdataConfig struct {
	ShowErro bool
	SafetheUrl bool
	SessionKey string
	SessionMaxAge int //时间分钟
	CookiePath string
	CookieDomain string
	CookieSameSite  http.SameSite
	CookieSecure bool
	CookieHttpOnly bool
	AllowUploadType []string
	AllowUploadMaxSize int //k 大小按照k计算
	UploadDir string
	MaxMultipartMemory int64

}
type RData struct {


	getData map[string]string
	postData map[string][]string
	cookieData []*http.Cookie
	Request *http.Request
	ResponseWriter http.ResponseWriter
	sessionStore  *sessions.CookieStore

	config rdataConfig

	tpl *tpl
	dbconn *dbConn
}

//处理数据
func (this *RData)parseQuery(){
	if this.getData != nil {
		return
	}

	//获取路由匹配上的参数
	theValue := mux.Vars(this.Request)

	//获取URL上传递的参数，会覆盖路由上的参数
	urlValues := this.Request.URL.Query()
	for k, v := range urlValues {

		if len(v) > 0 {
			theValue[k] = v[len(v)-1] //取传参的最后一个值
		} else {
			theValue[k] = ""
		}
	}

	//URL校验
	if this.config.SafetheUrl {
		for _,v := range theValue{
			match,err:=regexp.MatchString(`([\W]+)`,v)

			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
				}).Error("regexp.MatchString Error")
			}

			if match {
				log.WithFields(log.Fields{
					"Url":  this.Request.URL.String(),
				}).Info("Not Safe Url")
				theValue = nil
				//不安全的连接，应该统一到一个位置
				//panic("Url Not Safe")
				break
			}
		}
	}

	this.getData = theValue


}

func (this *RData) Query(v string) (string,bool) {

	this.parseQuery() //先处理数据

	if value, ok := this.getData[v]; ok{
		log.WithFields(log.Fields{
			"param":  v,
		}).Debug("Cant find params")

		return  value,ok
	}
	return "",false
}

func (this *RData)parseFormArray(){
	if this.postData != nil {
		return
	}
	this.Request.ParseForm()
	this.postData = this.Request.PostForm
}
func (this *RData) PostFormArray(v string) ([]string,bool) {
	this.parseFormArray()

	if valueList, ok := this.postData[v]; ok{
		log.WithFields(log.Fields{
			"param":  v,
		}).Debug("Cant find params")

		return  valueList,ok
	}
	return nil,false
}

func (this *RData) PostForm(v string) (string,bool) {

	if valueList, ok := this.PostFormArray(v); ok{
		log.WithFields(log.Fields{
			"param":  v,
		}).Debug("Cant find params")

		return  valueList[0],ok
	}
	return "",false
}


func (this *RData) Cookie(name string)(string,bool){
	value := ""
	find := false
	for _, c := range this.Request.Cookies() {
		if name == c.Name{
			value = c.Value
			find = true
		}
	}

	return value,find
}

// 简便方法设置cookie,通过配置文件里的配置
func (this *RData) SetCookie(name, value string,maxAge int){ //单位分钟
	path := this.config.CookiePath
	domain := this.config.CookieDomain
	sameSite := this.config.CookieSameSite
	secure := this.config.CookieSecure
	httpOnly := this.config.CookieHttpOnly
	this.SetCookieRaw(name,value,maxAge*60,path,domain,sameSite,secure,httpOnly)
}

// 基础方法， 注意这里的单位为秒
func (this *RData) SetCookieRaw(name, value string, maxAge int, path, domain string, sameSite http.SameSite, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(this.ResponseWriter, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (this *RData) initSession()  {
	if this.sessionStore == nil {
		SessionKey := this.config.SessionKey
		this.sessionStore = sessions.NewCookieStore([]byte(SessionKey))
		this.sessionStore.MaxAge(this.config.SessionMaxAge*60)
	}
}

func (this *RData)SetSession(key,value string) error  {
	this.initSession()

	session, _ := this.sessionStore.Get(this.Request, "sessionID")
	session.Values[key] = value
	err := session.Save(this.Request, this.ResponseWriter)
	return err
}

func (this *RData)Session(sessionKey string) string { //未获取到返回空字符串
	this.initSession()

	session, _ := this.sessionStore.Get(this.Request, "sessionID")
	//println(session.Values[sessionKey])
	if session.Values[sessionKey] == nil {
		return  ""
	}
	return session.Values[sessionKey].(string)
}

func (this *RData)ExecTpl(tplName string,maps map[string]interface{})  {
	this.tpl.Execute(this.ResponseWriter,tplName,maps)
}

var (
	UPLOAD_ERROR_NOFILE = errors.New("The File Not Find")
	UPLOAD_ERROR_UNKONW_FILE = errors.New("The File Is Unkonw")
	UPLOAD_ERROR_FILE_NOT_ALLOW = errors.New("The File Is Not Allow")
	UPLOAD_ERROR_FILE_BIG_THAN_MAX = errors.New("The File Is Big Than Allow")
	UPLOAD_ERROR_MKDIR_ERROR = errors.New("Make Dir error")
	UPLOAD_ERROR_SAVEFILE_ERROR = errors.New("Save file error")
)

func (this *RData) SimpleUpload(name string)(newFile,oldFileName string,fileSize int64,err error) {

	fh,fileErr := this.FormFile(name)

	//检查文件是否存在
	if fileErr != nil {
		err = UPLOAD_ERROR_NOFILE
		return
	}



	oldFileName = fh.Filename
	fileSize = fh.Size/1024

	//检查文件类型
	fileType,_ := this.GetFileType(fh)
	if fileType == "unknow" {
		err = UPLOAD_ERROR_UNKONW_FILE
		return
	}

	//检查文件类型是否合法
	isAllow := false
	for _,v := range this.config.AllowUploadType{
		if fileType == v {
			isAllow = true
			break
		}
	}

	if !isAllow {
		err = UPLOAD_ERROR_FILE_NOT_ALLOW
		return
	}

	//检查文件大小
	allowMaxSize :=	int64(this.config.AllowUploadMaxSize)*1024
	if allowMaxSize < fh.Size {
		err = UPLOAD_ERROR_FILE_BIG_THAN_MAX
		return
	}

	//存文件

	//先建目录
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir =  fmt.Sprintf("/%s/",Strip(dir,"/"))

	configDir := fmt.Sprintf("%s/",Strip(this.config.UploadDir,"/"))

	relativeDir := time.Now().Format("2006/01/02/")

	fullPath := dir+configDir+relativeDir
	err = os.MkdirAll( fullPath , os.ModePerm)
	if err != nil {
		err = UPLOAD_ERROR_MKDIR_ERROR
	}

	//得到文件名
	timeNowStr := convert.Int64ToString(time.Now().Unix())
	md5Byte := md5.Sum([]byte(timeNowStr))
	newFilename := fmt.Sprintf("%x.%s", md5Byte,fileType)

	err = this.SaveUploadedFile(fh, fullPath+newFilename)
	if err != nil {
		err = UPLOAD_ERROR_SAVEFILE_ERROR
		return
	}

	newFile =  relativeDir+newFilename

	return
}



//	没找到 文件 错误为http.ErrMissingFile
func (this *RData) FormFile(name string) (*multipart.FileHeader, error) {
	if this.Request.MultipartForm == nil {
		if err := this.Request.ParseMultipartForm(this.config.MaxMultipartMemory*1024); err != nil {
			return nil, err
		}
	}
	f, fh, err := this.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// 获取文件类型 没找到相关类型 返回 unknow
func (this *RData) GetFileType(fh *multipart.FileHeader) (fileType string,fileClass string){
	if fh == nil {
		panic("Get File Type Cant nil")
	}

	ContentType := fh.Header.Get("Content-Type")
	mimeInfo := mime.GetFileTypeByMIME(ContentType)

	fileType = mimeInfo[0]
	fileClass = mimeInfo[1]
	return
}

// 保存文件
func (this *RData) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	fmt.Println("\n"+dst)
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
/*
//ishas,newfile,oldfile,err := upfile.UploadFile("kkk",r)
//fmt.Printf("\n%v,%v,%v,%v\n",ishas,newfile,oldfile,err)
func UploadFile(upKey string,r *http.Request) (isHas bool,NewFile string,Oldname string,err error) {//是否上传，新的文件名，老的文件名，错误信息

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir += "/"

	configDir := cfg.GetString("UploadDir")+"/"

	relativeDir := time.Now().Format("2006/01/02/")
	err = os.MkdirAll( dir+configDir+relativeDir , os.ModePerm)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"filePath":  dir+configDir+relativeDir,
		}).Error("Mkdir error")
		return true,"","",err
	}
	timeNowStr := convert.Int64ToString(time.Now().Unix())
	md5Byte := md5.Sum([]byte(timeNowStr))
	newFilename := fmt.Sprintf("%x", md5Byte)
	//


	fileObject := new(Upload)

	fileObject.UpConfig = new(Upconfig)
	fileObject.UpConfig.AllowType = cfg.GetSliceString("AllowUploadType")
	fileObject.UpConfig.AllowMaxSize = cfg.GetInt("AllowUploadSize")
	//应该改为返回状态码
	isHas,fileExt,oldname,err := fileObject.Save(dir+configDir+relativeDir,newFilename,upKey,r)
	if err != nil && isHas {
		log.WithFields(log.Fields{
			"err":  err,
			"filePath":  dir+configDir+relativeDir+newFilename+"."+fileExt, //
		}).Debug("Upload file error")
		return isHas,"",oldname,err
	}


	return isHas,relativeDir+newFilename+"."+fileExt,oldname,nil
}






func (this *Upload)Save(dir,newFilename,upKey string,r *http.Request) (isHas bool,newType string,oldName string,err error)  {

	mfile,partHead,err := r.FormFile(upKey)
	if err != nil {
		isHas = false
		return
	}else {
		isHas = true
	}

	defer mfile.Close()

	//fileName := partHead.Filename
	size := partHead.Size


	contentType := partHead.Header["Content-Type"]

	//不在允许的类型范围内
	fileType,isAllow := this.isInAllow(contentType)
	if !isAllow {
		err = errors.New("The FileType is not allow")
		return
	}

	//超过大小
	if int(math.Ceil(float64(size)/float64(1024.0))) > this.UpConfig.AllowMaxSize {
		err = errors.New("The FileType is Max Than"+string(this.UpConfig.AllowMaxSize)+"k")
		return
	}

	f, errOpen := os.OpenFile(dir + newFilename+"."+fileType, os.O_WRONLY | os.O_CREATE, 0666)
	if errOpen != nil{
		err = errOpen
		return
	}
	defer f.Close()
	//拷贝文件
	io.Copy(f, mfile)

	return
}

func (this *Upload)isInAllow(contentType []string) (string,bool) {

	allow := true
	fileType := ""
	for _,value := range contentType{//contentType竟然是个数组
		ctype := GetFileTypeByMIME(value)
		inAllow := false
		for _,v := range this.UpConfig.AllowType{
			if ctype[0] == v {
				fileType = ctype[0]
				inAllow = true
				break
			}
		}
		if inAllow == false{
			allow = false
			break
		}
	}

	return fileType,allow

}*/




