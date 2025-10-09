package file

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileFinder struct {
	Root string
}
type StatefulMap struct {
	mu   sync.Mutex
	Data map[string]*FileStatus
}

type FileStatus struct {
	Path    string
	Started bool
	Done    bool
}

func (sm *StatefulMap) GetNext() (string, *FileStatus, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for filename, fs := range sm.Data {
		if !fs.Started {
			fs.Started = true
			return filename, fs, true
		}
	}
	return "", nil, false
}

func (f *FileFinder) FindByExtension(extension string) StatefulMap {
	files := make(map[string]*FileStatus)
	filepath.WalkDir(f.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return filepath.SkipDir
			}
			return err
		}

		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		if !d.IsDir() && filepath.Ext(d.Name()) == extension {
			files[d.Name()] = &FileStatus{
				Path:    path,
				Started: false,
				Done:    false,
			}
		}
		return nil
	})
	return StatefulMap{
		Data: files,
	}
}
