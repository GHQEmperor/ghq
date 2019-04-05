package ghq

type StaticFile struct {
	// uri path.
	Uri     string

	// local file path.
	DirPath string
}

// set static file path.
// all files of this uri is a public file.
// demo: g.SetStaticFile("/static/","static")
func (r *Router) SetStaticFile(uri, dirPath string) {
	r.staticFileUri = append(r.staticFileUri, StaticFile{Uri:uri,DirPath:dirPath})
}