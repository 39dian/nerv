package model

import (
	"path/filepath"
	"os"
	"strings"
	"encoding/json"
	"github.com/toolkits/file"
)

// ClassRepository load classes from file path
type ClassRepository struct {
	path    string
	classes map[string]*Class
}

var classRep *ClassRepository

//InitClassRepository initialize ClassRepository and setup the path of classes
func InitClassRepository(path string) {
	if path == "" {
		panic("path must not be nil")
	}
	classRep = &ClassRepository{path, map[string]*Class{}}

	filepath.Walk(path, loadClassFiles)
}

func loadClassFiles(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.HasSuffix(path, ".json") {
		content, e := file.ToTrimString(path)
		if e != nil {
			return e
		}

		var class Class
		e = json.Unmarshal([]byte(content), &class)
		if e != nil {
			return e
		}
		classRep.classes[class.Name] = &class
		return nil
	} else {
		return nil
	}
}

//ClassRepository return the singleton
func GetClassRepository() *ClassRepository {
	return classRep
}

//GetClass return type of the name or nil if isn't  exists.
func (p *ClassRepository) Find(name string) *Class {
	return p.classes[name]
}
