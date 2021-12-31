package web

import (
	"image/jpeg"
	"log"
	"net/http"
	"strings"
)

func (s Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	stdFolders := s.folders.GetStandardFolders()

	s.tplExec("index.html", w, map[string]interface{}{
		"Folders": stdFolders,
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

	m, err := folder.GetMedia(mediaName)
	if err != nil {
		log.Printf("Media get for '%s': %v", folder.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i, err := m.Thumbnail()
	if err != nil {
		log.Printf("Thumbnail for '%s': %v", m.Title, err)
		return
	}

	jpeg.Encode(w, i, nil)
}

func (s Server) handleFullPicture(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	mediaName := strings.TrimPrefix(r.URL.Path, "/images/")
	mediaName = strings.TrimSuffix(mediaName, ".mp4") // As a hack, we add ".mp4" to the end of all videos to trick the gallery into displaying them as videos

	parts := strings.Split(mediaName, "/")

	folder, ok := s.folders.Get(parts[0])
	if !ok {
		log.Printf("Requested unknown folder: %s", r.FormValue("f"))
		s.tplExec("notfound.html", w, nil)
		return
	}

	mediaName = parts[1]
	m, err := folder.GetMedia(mediaName)
	if err != nil {
		log.Printf("Media get for '%s': %v", folder.Name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = m.BrowserSafeMedia(w)
	if err != nil {
		log.Printf("BrowserSafeMedia for '%s' in '%s': %v", mediaName, folder.Name, err)
	}
}
