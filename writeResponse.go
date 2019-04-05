package ghq

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
)

// Html put in views dir.
// write Html.
// rw.WriteHtml("index.html").
func (rw *RW) WriteHTML(fileName string) (err error) {
	file ,err := os.Open("views/"+fileName)
	if err != nil {
		return
	}
	defer file.Close()
	fileBytes ,err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	rw.W.WriteHeader(200)
	_,err = rw.W.Write(fileBytes)
	return
}

// write Json.
func (rw *RW) WriteJson(data interface{}) (err error) {
	jsonBytes,err := json.Marshal(data)
	if err != nil {
		return
	}
	rw.W.WriteHeader(200)
	_,err = rw.W.Write(jsonBytes)
	return
}

// write XML.
func (rw *RW) WriteXML(data interface{}) (err error) {
	xmlBytes,err := xml.Marshal(data)
	if err != nil {
		return
	}
	rw.W.WriteHeader(200)
	_,err = rw.W.Write(xmlBytes)
	return
}