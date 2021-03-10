package filesystem

import (
	"embed"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/macaron.v1"
)

var _ macaron.TemplateFileSystem = (*FS)(nil)

type FS struct {
	embed embed.FS
}

func NewFS(embedFS embed.FS) *FS {
	return &FS{
		embed: embedFS,
	}
}

func (fs *FS) ListFiles() []macaron.TemplateFile {
	return fs.scanDir(".")
}

func (fs *FS) scanDir(dirPath string) []macaron.TemplateFile {
	templateFiles := make([]macaron.TemplateFile, 0)

	dirEntry, err := fs.embed.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	for _, entry := range dirEntry {
		entryPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() {
			templateFiles = append(templateFiles, fs.scanDir(entryPath)...)
			continue
		}

		fileReader, err := fs.embed.Open(entryPath)
		if err != nil {
			continue
		}
		fileData, err := io.ReadAll(fileReader)
		if err != nil {
			continue
		}

		f := &file{
			DirEntry:  entry,
			entryPath: entryPath,
			data:      fileData,
		}
		templateFiles = append(templateFiles, f)
	}
	return templateFiles
}

func (fs *FS) Get(fileName string) (io.Reader, error) {
	return fs.embed.Open(fileName)
}

var _ macaron.TemplateFile = (*file)(nil)

type file struct {
	fs.DirEntry
	entryPath string
	data      []byte
}

func (f *file) Name() string {
	if f.Ext() == ".html" || f.Ext() == ".tmpl" {
		return strings.TrimSuffix(f.entryPath, f.Ext())
	}
	return f.entryPath
}

func (f *file) Data() []byte {
	return f.data
}

func (f *file) Ext() string {
	return filepath.Ext(f.DirEntry.Name())
}
