package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	/*
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.save()
		p2, _ := loadPage("TestPage")
		fmt.Println(string(p2.Body))
	*/
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		//http.Redirect(w, r, "/view/"+title, http.StatusFound)
		//fmt.Println("landed on %s", title)
		p = &Page{Title: title, Body: []byte("status")}
		renderAppTemplate(w, "404", p)
		return
	}

	//hardcoded html code
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)

	//templator
	renderAppTemplate(w, "view", p)

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Println("Body = %s", p.Body)
	//log.Fatalln("Body = %s", p.Body)

	//hardcoded html code
	/*fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	"<form action=\"/save/%s\" method=\"POST\">"+
	"<textarea name=\"body\">%s  1</textarea><br>"+
	"<input type=\"submit\" value=\"Save\">"+
	"</form>",
	p.Title, p.Title, p.Body)
	*/

	//templator
	renderAppTemplate(w, "edit", p)

}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderAppTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("templates/" + tmpl + ".html")
	t.Execute(w, p)
}
