package rest

import(
	"embed"
	"io"
	"io/fs"
	"fmt"
	S "strings"
	)

//go:embed embedfs
var EmbedFS embed.FS
// data, _ := f.ReadFile("hello.txt")

func init() {
     ListFS(&EmbedFS)
     }

/*
https://pkg.go.dev/embed#hdr-File_Systems

FS implements interface io/fs.FS, so it can be used
with any package that understands file systems,
including net/http, text/template, html/template.

func (f FS) Open(name string) (fs.File, error)
   (The returned file implements io.Seeker and
     io.ReaderAt when the file is not a directory)
func (f FS) ReadDir(name string) ([]fs.DirEntry, error)
func (f FS) ReadFile(name string) ([]byte, error)

For example,
http.Handle(
	"/static/",
	http.StripPrefix("/static/",
	http.FileServer(http.FS(content))))
Or,
	content1, _ := folder.ReadFile("folder/file1.hash")
	print(string(content1))
*/

func ListFS(pFS *embed.FS) error {
     var DE []fs.DirEntry
     var e error
     
     DE, e = pFS.ReadDir(".")
     if e != nil {
     	return fmt.Errorf("ListFS(.): %w", e)
	}
     fmt.Printf("dot: %+v \n", DE)
     
     DE, e = pFS.ReadDir("embedfs/h-body")
     if e != nil {
     	return fmt.Errorf("ListFS(h-body): %w", e)
	}
     fmt.Printf("h-body: \n%+v \n", DE)
     
     DE, e = pFS.ReadDir("embedfs/h-frag")
     if e != nil {
     	return fmt.Errorf("ListFS(h-frag): %w", e)
	}
     fmt.Printf("h-frag: \n%+v \n", DE)
     return nil
     }

func GetBody(s string) (string, error) {
     if !S.HasSuffix(s, ".hb") {
     	s += ".hb"
	}
     f,e := EmbedFS.Open("embedfs/h-body/" + s)
     if e != nil {
     	return "", fmt.Errorf("rest.embed.GetBody(%s): %w", s, e)
	}
     bb, ee := io.ReadAll(f)
     return string(bb), ee
     }

func GetFrag(s string) (string, error) {
     if !S.HasSuffix(s, ".hf") {
     	s += ".hf"
	}
     f,e := EmbedFS.Open("embedfs/h-frag/" + s)
     if e != nil {
     	return "", fmt.Errorf("rest.embed.GetFrag(%s): %w", s, e)
	}
     bb, ee := io.ReadAll(f)
     return string(bb), ee
     }

