package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// Use ParseGlob to parse all templates, including partials
//var templates = template.Must(template.ParseGlob("templates/**/*.html"))

var templates = template.Must(template.ParseFiles(
	"templates/base.html",
	"templates/home.html",
	"templates/partials/footer.html",
	// Add other template paths here if needed
))

type Page struct {
	Title     string
	Body      []byte
	Timestamp int64
}

func main() {
	/*
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.save()
		p2, _ := loadPage("TestPage")
		fmt.Println(string(p2.Body))
	*/
	http.HandleFunc("/", defaultViewHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body, Timestamp: time.Now().Unix()}, nil
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func defaultViewHandler(w http.ResponseWriter, r *http.Request) {
	results := []string{}
	filesWiki, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, fileWiki := range filesWiki {
		fileWikiNameLen := len(fileWiki.Name()) - 4
		if fileWiki.Name()[fileWikiNameLen:] == ".txt" {
			fmt.Println(fileWiki.Name())
			results = append(results, fileWiki.Name()[:fileWikiNameLen])
		}
	}

	m := map[string]interface{}{
		"Results":   results,
		"Title":     "Home Page",
		"Timestamp": time.Now().Unix(),
	}

	//t, _ := template.ParseFiles("templates/default.html")
	//t.Execute(w, m)
	//fmt.Println(m)

	//p := &Page{Title: "Test Page", Body: []byte("This is a test page body")}
	renderAppTemplate(w, "home", m)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		//http.Redirect(w, r, "/view/"+title, http.StatusFound)
		//fmt.Println("landed on %s", title)
		p = &Page{Title: title, Body: nil, Timestamp: time.Now().Unix()}
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
		p = &Page{Title: title, Timestamp: time.Now().Unix()}
	}

	//fmt.Println("Body = %s", p.Body)
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

func renderAppTemplate(w http.ResponseWriter, template string, data interface{}) {
	//the following appraoch can't load 'partials' so it's replaced with `ExecuteTemplate`
	//t, _ := template.ParseFiles("templates/" + tmpl + ".html")
	//t.Execute(w, p)

	// Parse the base template
	/*
		tmpl, err := template.New("base.html").ParseFiles("templates/base.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	tmpl, _ := templates.Clone()

	// Parse the content (child) template if specified
	if template != "" {
		_, err := tmpl.ParseFiles("templates/" + template + ".html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//add footer
	_, err := tmpl.ParseFiles("templates/partials/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the base template
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
