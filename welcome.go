package snailframe

func Welcome(r *RData)  {

	m,_ := r.Query("aa")
	//re,_ := r.dbconn.Find("SELECT * FROM shici_info where id=?",m)
	//fmt.Println(re)

	re2 := map[string]interface{}{
		"a":m,
	}
	r.ExecTpl("aaa",re2)


}
