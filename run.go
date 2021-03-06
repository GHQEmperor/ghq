package ghq

import (
	"errors"
	"fmt"
	"net/http"
)

func New() *Router {
	var RouterC Router
	//RouterC.Config = make(map[string]string)
	if err := RouterC.LoadConfig(); err != nil {
		panic(err)
	}
	RouterC.functionsM = make(funcMap)
	RouterC.uriFuncM = make(map[string]httpRW)
	return &RouterC
}

// run the http server.
func (r *Router) Run() error {
	// set the uri function.
	for uri, fn := range r.functionsM {
		r.newUri(uri, fn)
	}
	// register the uri function  to http package.
	for uri, uriFn := range r.uriFuncM {
		http.HandleFunc(uri, uriFn)
	}

	// register the static file path to http package.
	for _, sf := range r.staticFileUri {
		http.Handle(sf.Uri, http.StripPrefix(sf.Uri, http.FileServer(http.Dir(sf.DirPath))))
	}
	// TODO: custom 404 page

	// set serve port.
	port, ok := GetConfig("port")
	if !ok {
		return errors.New("don't appoint port")
	}
	fmt.Println("run port ", port)
	return http.ListenAndServe(port, nil)
}
