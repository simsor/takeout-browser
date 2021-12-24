package web

import (
	"embed"
	"html/template"
	"image/jpeg"
	"log"
	"net/http"
	"strings"

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

func (s Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.tplExec("index.html", w, map[string]interface{}{
		"Folders": s.folders,
	})
}

func (s Server) handleFolder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	folder, ok := s.folders.Get(r.FormValue("f"))
	if !ok {
		s.tplExec("notfound.html", w, nil)
		return
	}

	s.tplExec("folder.html", w, map[string]interface{}{"Folder": &folder})
}

func (s Server) handleThumbnail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mediaName := strings.TrimPrefix(r.URL.Path, "/thumbs/")
	folder, ok := s.folders.Get(r.FormValue("f"))
	if !ok {
		log.Printf("Requested unknown folder: %s", r.FormValue("f"))
		s.tplExec("notfound.html", w, nil)
		return
	}

	media, err := folder.Media()
	if err != nil {
		log.Printf("Media get for '%s': %v", folder.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, m := range media {
		if m.Title == mediaName {
			i, err := m.Thumbnail()
			if err != nil {
				log.Printf("Thumbnail for '%s': %v", m.Title, err)
				return
			}

			jpeg.Encode(w, i, nil)
			return
		}
	}

	log.Printf("Not found: '%s' in '%s'", mediaName, folder.Name)
}

func (s Server) handleFullPicture(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mediaName := strings.TrimPrefix(r.URL.Path, "/images/")
	folder, ok := s.folders.Get(r.FormValue("f"))
	if !ok {
		log.Printf("Requested unknown folder: %s", r.FormValue("f"))
		s.tplExec("notfound.html", w, nil)
		return
	}

	media, err := folder.Media()
	if err != nil {
		log.Printf("Media get for '%s': %v", folder.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, m := range media {
		if m.Title == mediaName {
			err = m.BrowserSafeMedia(w)
			if err != nil {
				log.Printf("BrowserSafeMedia for '%s' in '%s': %v", mediaName, folder.Name, err)
			}
			return
		}
	}

	log.Printf("Not found: '%s' in '%s'", mediaName, folder.Name)
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
