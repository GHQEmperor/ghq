package ghq

import "net/http"

type Router struct {
	// configure information.
	Config map[string]string

	// storage get,post,head,put... function.
	// a uri correspond to a function slice.
	// the length of function slice is MethodCount.
	functionsM funcMap

	// storage uri function.
	// this uri function consist of a functionsM slice.
	uriFuncM      map[string]httpRW
	staticFileUri []StaticFile
}

type RW struct {
	W           http.ResponseWriter
	R           *http.Request
	isParseForm bool
	*Router
}

type FuncRW func(rw RW)
type httpRW func(w http.ResponseWriter, r *http.Request)
type funcMap map[string]*[]FuncRW

const (
	MethodGet = iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
)

const MethodCount = 9

// if the slice is not initialized,this function initializes the slice.
func checkMethodInit(uri string, f funcMap) {
	_, ok := f[uri]
	if !ok {
		f[uri] = newFuncRWs()
	}
}

func (r *Router) Get(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodGet] = function
}

func (r *Router) Head(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodHead] = function
}

func (r *Router) Post(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodPost] = function
}

func (r *Router) Put(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodPut] = function
}

func (r *Router) Patch(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodPatch] = function
}

func (r *Router) Delete(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodDelete] = function
}

func (r *Router) Connect(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodConnect] = function
}

func (r *Router) Options(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodOptions] = function
}

func (r *Router) Trace(uri string, function FuncRW) {
	checkMethodInit(uri, r.functionsM)
	(*r.functionsM[uri])[MethodOptions] = function
}

func newFuncRWs() *[]FuncRW {
	funcRw := make([]FuncRW, MethodCount)
	for i := 0; i < MethodCount; i++ {
		funcRw[i] = _404
	}
	return &funcRw
}

// create a function: func(w http.ResponseWriter, r *http.Request)
// to register http.HandleFunc.
func (r *Router) newUri(uri string, functions *[]FuncRW) {
	r.uriFuncM[uri] = func(responseWriter http.ResponseWriter, request *http.Request) {
		rw := RW{responseWriter, request, false, r}
		switch request.Method {
		case http.MethodGet:
			(*functions)[MethodGet](rw)
		case http.MethodHead:
			(*functions)[MethodHead](rw)
		case http.MethodPost:
			(*functions)[MethodPost](rw)
		case http.MethodPut:
			(*functions)[MethodPut](rw)
		case http.MethodPatch:
			(*functions)[MethodPatch](rw)
		case http.MethodDelete:
			(*functions)[MethodDelete](rw)
		case http.MethodConnect:
			(*functions)[MethodConnect](rw)
		case http.MethodOptions:
			(*functions)[MethodOptions](rw)
		case http.MethodTrace:
			(*functions)[MethodTrace](rw)
		default:
			_404(rw)
		}
	}
}

func _404(rw RW) {
	rw.W.WriteHeader(404)
	_, _ = rw.W.Write([]byte("404"))
}
