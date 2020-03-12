package snailframe

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"strings"
)

type dbConn struct {
	conn *sql.DB
}

type dbConfig struct {
	UserName string
	Password string
	Host string
	Port int64
	DataBase string
	CharSet string
}

// 初始化数据库，连接数
func newDB(cfg dbConfig) *dbConn {

	UserName := cfg.UserName
	Password := cfg.Password
	Host := cfg.Host
	Port := cfg.Port
	DataBase := cfg.DataBase
	CharSet := cfg.CharSet

	log.WithFields(log.Fields{
		"dbConfig":     cfg,
	}).Trace("Start Connet Database")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", UserName, Password,
		Host, Port, DataBase, CharSet)

	conect, err := sql.Open("mysql", dataSource)
	conect.SetMaxOpenConns(0)
	conect.SetMaxIdleConns(1000)//https://blog.csdn.net/wangguoyang429883793/article/details/73436563/
	if err != nil {
		log.WithFields(log.Fields{
			"Host":     Host,
			"Port":     Port,
			"DataBase": DataBase,
		}).Panic("Connet Database Falt")

		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	if conect.Ping() != nil {

		log.WithFields(log.Fields{
			"Host":     Host,
			"Port":     Port,
			"DataBase": DataBase,
		}).Panic("Connet Database Falt")

		panic(err.Error())
	}else{
		log.WithFields(log.Fields{
			"Host":     Host,
			"Port":     Port,
			"DataBase": DataBase,
		}).Trace("Connet Database success")
	}



	return &dbConn{conect}
}

/**
使用方法
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "db"

					  账号:密码@tcp(数据库地址:端口号)/数据库名称
db2, err := sql.Open("mysql", "root:502884@tcp(127.0.0.1:3306)/zdm")
if err != nil {
	panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
}
mysqldb_obj  :=    new(db.mysqldb);
mysqldb_obj.Conn = db2;
b,err :=  mysqldb_obj.GetFirstField("SELECT  uid  FROM user where uid=2 order by uid asc limit 5");
if err != nil {
	panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
}
**/

//返回整个查询结果
func (this *dbConn)FindAll(query_sql string, args ...interface{}) (result []map[string]string, err error) {
	//obj.mux.Lock()
	//defer obj.mux.Unlock()

	log.WithFields(log.Fields{
		"SQL":  query_sql,
		"Args": args,
	}).Debug("Start SQL FindAll")

	// Execute the query
	rows, err := this.conn.Query(query_sql, args...)
	defer rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  query_sql,
			"Args": args,
			"err":  err,
		}).Error("Query sql error")

		return result, err
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows

	// var  result_key = 0;
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, err
		}
		var value string
		var row map[string]string = make(map[string]string)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			row[columns[i]] = value
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		return result, err // proper error handling instead of panic in your app
	}
	return result, nil
}

//只返回第一条查询结果
func (this *dbConn)Find(query_sql string, args ...interface{}) (result map[string]string, err error) {
	//obj.mux.Lock()
	//defer obj.mux.Unlock()
	// Execute the query

	log.WithFields(log.Fields{
		"SQL":  query_sql,
		"Args": args,
	}).Debug("Start SQL Find")

	rows, err := this.conn.Query(query_sql, args...)
	defer rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  query_sql,
			"Args": args,
			"err":  err,
		}).Panic("Query sql error")
		return result, err
	}
	// Get column names
	columns, err := rows.Columns()

	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  query_sql,
			"Args": args,
			"err":  err,
		}).Error("Query sql error")
		return result, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows

	var row map[string]string = make(map[string]string)
	// var  result_key = 0;
	defer  rows.Close()
	if rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, err
		}
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			row[columns[i]] = value
		}
		return row, nil
	}
	return row, nil
}

//只返回第一条查询结果第一个字段
func (this *dbConn)GetFirstField(query_sql string, args ...interface{}) (result string, err error) {
	//obj.mux.Lock()
	//defer obj.mux.Unlock()
	// Execute the query

	log.WithFields(log.Fields{
		"SQL":  query_sql,
		"Args": args,
	}).Debug("Start SQL GetFirstField")

	rows, err := this.conn.Query(query_sql, args...)
	defer rows.Close()
	if err != nil {
		return result, err
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  query_sql,
			"Args": args,
			"err":  err,
		}).Error("Query sql error")
		return result, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows

	var value string
	// var  result_key = 0;
	if rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, err
		}
		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			return value, nil
		}
	}
	return value, nil
}

func (this *dbConn)Exec(querySql string, args ...interface{}) (RowsAffected int64, err error) {
	//obj.mux.Lock()
	//defer obj.mux.Unlock()

	log.WithFields(log.Fields{
		"SQL":  querySql,
		"Args": args,
	}).Debug("Start SQL Exec")

	var LastInsertId int64
	LastInsertId = 0
	rs, err := this.conn.Exec(querySql, args...)
	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  querySql,
			"Args": args,
			"err":  err,
		}).Error("Query sql error")

		return 0, err
	}
	c, err := rs.LastInsertId()
	d, err := rs.RowsAffected()
	if c > 0 {
		LastInsertId = c

		log.WithFields(log.Fields{
			"SQL":          querySql,
			"Args":         args,
			"LastInsertId": LastInsertId,
		}).Debug("Insert SQL")

	}
	if LastInsertId == 0 && d > 0 {

		log.WithFields(log.Fields{
			"SQL":          querySql,
			"Args":         args,
			"RowsAffected": LastInsertId,
		}).Debug("Exec SQL")

		LastInsertId = d
	}

	return LastInsertId, nil
}

func (this *dbConn)Insert(table string, data map[string]interface{}) (LastInsertId int64, err error) {
	//obj.mux.Lock()
	//defer obj.mux.Unlock()
	LastInsertId = 0
	var fieldArr []string
	var holderArr []string
	var valueArr []interface{}
	for k, v := range data {
		fieldArr = append(fieldArr, "`"+k+"`")
		holderArr = append(holderArr, "?")
		valueArr = append(valueArr, v)
	}
	var _sql = "INSERT INTO  `" + table + "`   (" + strings.Join(fieldArr, ",") + ") VALUES(" + strings.Join(holderArr, ",") + ")"

	log.WithFields(log.Fields{
		"SQL":  _sql,
		"Args": valueArr,
	}).Debug("Insert query sql")

	rs, err := this.conn.Exec(_sql, valueArr...)

	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  _sql,
			"Args": valueArr,
			"err":  err,
		}).Error("Query sql error")

		return 0, err
	}
	// LastInsertId() (int64, error)

	// RowsAffected返回被update、insert或delete命令影响的行数。
	// 不是所有的数据库都支持该功能。
	// RowsAffected() (int64, error)
	c, err := rs.LastInsertId()
	d, err := rs.RowsAffected()
	if c > 0 {
		LastInsertId = c
	}
	if LastInsertId == 0 && d > 0 {
		LastInsertId = d
	}
	return LastInsertId, nil
}

func (this *dbConn)UpdateByWhereMap(table string, data map[string]interface{}, whereMap map[string]interface{}) (RowsAffected int64, err error) {

	var updateString = ""

	var valueArr []interface{}

	for k, v := range data {
		if updateString == "" {
			updateString += "`" + k + "`=?"
		} else {
			updateString += " ,`" + k + "`=?"
		}
		valueArr = append(valueArr, v)
	}

	var WhereString = "1"

	for k, v := range whereMap {
		WhereString += " and `" + k + "`= ?"
		valueArr = append(valueArr, v)
	}

	var querySql string = "UPDATE  " + table + " SET " + updateString + "  WHERE " + WhereString

	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  querySql,
			"Args": valueArr,
			"err":  err,
		}).Debug("Update sql")
	}

	RowsAffected, err = this.Exec(querySql, valueArr...)

	if err != nil {
		log.WithFields(log.Fields{
			"SQL":  querySql,
			"Args": valueArr,
			"err":  err,
		}).Error("Query sql error")
	}
	return RowsAffected, err
}

// 多段sql，拼接 a.sql("sss",sss).sql(" AND ",ssss).Execute()
type mutiSql struct {
	sql  string
	args []interface{}
	con *dbConn
}

func (this *dbConn)NewMutiSql() *mutiSql {
	muti :=  new(mutiSql)
	muti.con = this
	return muti
}

func (this *mutiSql) SQL(sql string, args ...interface{}) *mutiSql {
	if len(sql) > 0 {
		this.sql += " "
	}

	this.sql += sql

	this.args = append(this.args, args...)
	return this
}

func (this *mutiSql) Exec() (RowsAffected int64, err error) {
	return this.con.Exec(this.sql, this.args...)
}

func (this *mutiSql) GetFirstField() (result string, err error) {
	return this.con.GetFirstField(this.sql, this.args...)
}

func (this *mutiSql) Find() (result map[string]string, err error) {
	return this.con.Find(this.sql, this.args...)
}

func (this *mutiSql) FindAll() (result []map[string]string, err error) {
	return this.con.FindAll(this.sql, this.args...)
}

