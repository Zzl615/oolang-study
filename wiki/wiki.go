/**
 * @Author: Noaghzil
 * @Date:   2022-09-25 15:32:17
 * @Last Modified by:   Noaghzil
 * @Last Modified time: 2022-09-25 23:17:27
 */
// 读写文件
package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	// The title is the second subexpression.
	return m[2], nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, %s!", r.URL.Path[1:])
}

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, err := template.ParseFiles(tmpl + ".html")
// 	if err != nil {
// 		// 错误处理： http error status code
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = t.Execute(w, p)
// 	if err != nil {
// 		// 错误处理： http error status code
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	// ParseFiles更好的方法是ParseFiles在程序初始化时调用一次
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		// 错误处理： http error status code
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		// 错误处理：http跳转
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		// 错误处理：没有内容
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":6150", nil))
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()

}
