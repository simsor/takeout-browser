package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/simsor/takeout-browser/takeout"
)

var (
	//go:embed static/*
	static embed.FS
	//go:embed templates/*
	templates embed.FS
)

type Server struct {
	tpl     *template.Template
	folders takeout.Folders
}

func NewServer(folders takeout.Folders) *Server {
	s := Server{
		folders: folders,
	}
	s.loadTemplates()

	return &s
}

func (s Server) ListenAndServe(listen string) error {
	fileServer := http.FileServer(http.FS(static))

	http.Handle("/static/", fileServer)
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/folder", s.handleFolder)
	http.HandleFunc("/thumbs/", s.handleThumbnail)
	http.HandleFunc("/images/", s.handleFullPicture)

	return http.ListenAndServe(listen, nil)
}

func (s *Server) loadTemplates() {
	s.tpl = template.Must(template.ParseFS(templates, "templates/*.html"))
}

func (s Server) tplExec(name string, w http.ResponseWriter, params map[string]interface{}) {
	tpl := s.tpl.Lookup(name)
	if tpl == nil {
		log.Printf("Could not find template '%s'", name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := tpl.Execute(w, params)
	if err != nil {
		log.Printf("Error executing '%s': %v", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
