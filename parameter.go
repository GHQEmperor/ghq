package ghq

import (
	"errors"
	"strconv"
)

func (rw *RW) GetString(key string) (value string) {
	rw.isparseForm()
	value = rw.R.FormValue(key)
	if value == "" {
		value = rw.R.PostFormValue(key)
	}
	return
}

func (rw *RW) GetInt(key string) (value int, err error) {
	//rw.isparseForm()
	//valueStr := rw.R.FormValue(key)
	//if valueStr == "" {
	//	valueStr = rw.R.PostFormValue(key)
	//	if valueStr == "" {
	//		return 0, errors.New("no this int parameter")
	//	}
	//}
	valueStr := rw.GetString(key)
	if valueStr == "" {
		return 0, errors.New("no this int parameter")
	}
	// TODO:if Atoi failed ,the server will panic?
	return strconv.Atoi(valueStr)
}

func (rw *RW) isparseForm() {
	if !rw.isParseForm {
		_ = rw.R.ParseForm()
		rw.isParseForm = true
	}
}
