package takeout

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// Folder represents a Google Photo Takeout folder. It contains Media
type Folder struct {
	Name string

	path       string
	mediaCache []Media
}

type Folders []*Folder

// LoadPhotoTakeoutFolders loads all Google Photos folders in the given Takeout folder
func LoadPhotoTakeoutFolders(path string) (Folders, error) {
	var folders []*Folder

	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat: %v", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", path)
	}

	_, err = os.Stat(filepath.Join(path, "archive_browser.html"))
	if err != nil {
		return nil, fmt.Errorf("os.Stat: %v", err)
	}

	photosDir := filepath.Join(path, "Google\u00a0Photos")
	stat, err = os.Stat(photosDir)
	if err != nil {
		return nil, fmt.Errorf("os.Stat: %v", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", photosDir)
	}

	tkFolders, err := os.ReadDir(photosDir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	for _, entry := range tkFolders {
		if !entry.IsDir() {
			continue
		}

		path := filepath.Join(photosDir, entry.Name())
		f := Folder{
			path: path,
			Name: entry.Name(),
		}

		folders = append(folders, &f)
	}

	return Folders(folders), nil
}

// Media returns all media in the folder, ordered by date (decreasing)
func (f *Folder) Media() ([]Media, error) {
	if f.mediaCache != nil {
		return f.mediaCache, nil
	}

	var medias []Media

	glob, err := filepath.Glob(filepath.Join(f.path, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("filepath.Glob: %v", err)
	}

	for _, m := range glob {
		if filepath.Base(m) == "metadata.json" || filepath.Base(m) == "métadonnées.json" {
			continue
		}

		media, err := OpenMedia(m)
		if err != nil {
			log.Printf("open media '%s': %v", m, err)
			continue
		}

		medias = append(medias, media)
	}

	sort.Sort(ByReverseTimestamp(medias))

	f.mediaCache = medias
	return medias, nil
}

// GetMedia looks for a media in the folder with the given name
func (f *Folder) GetMedia(name string) (Media, error) {
	media, err := f.Media()
	if err != nil {
		return Media{}, err
	}

	for _, m := range media {
		if m.Title == name {
			return m, nil
		}
	}

	return Media{}, fmt.Errorf("could not find anything named '%s'", name)
}

// IsStandardPictureFolder checks whether this folder is from the main Google Photos folder
func (f Folder) IsStandardPictureFolder() bool {
	r := regexp.MustCompile(`^Photos from [1-2][0-9]{3}$`)
	return r.MatchString(f.Name)
}

// Get returns the folder with the given name. The boolean will be True if the folder was found.
func (fs Folders) Get(name string) (*Folder, bool) {
	for _, f := range fs {
		if f.Name == name {
			return f, true
		}
	}

	return nil, false
}

// GetStandardFolders returns all folders part of the main Google Photos folder
func (fs Folders) GetStandardFolders() []*Folder {
	var s []*Folder

	for _, f := range fs {
		if f.IsStandardPictureFolder() {
			s = append(s, f)
		}
	}

	// We want to return the folders in the reverse order: this way, when iterating over all folders
	// and all medias inside it, it will be ordered from most recent to oldest

	var rev []*Folder
	for i := len(s) - 1; i >= 0; i-- {
		rev = append(rev, s[i])
	}

	return rev
}
