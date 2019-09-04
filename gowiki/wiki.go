package main 


import (
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "html/template"
        "regexp"
        "errors"
)




//data structure that need to be implemented
type  Page struct {
        Title string 
        Body []byte
}

//global variable and initialize it with ParserFiles
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//globla to validate regex path 
var validPath = regexp.MustCompile("^(edit|save|view)/([a-zA-Z0-9]+)$")

//save the page on disk
func (p *Page) save() error {
        filename := p.Title + ".txt"
        return  ioutil.WriteFile(filename, p.Body, 0600)
}






//get the title parsing the path 
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
                http.NotFound(w,r)
                return "", errors.New("Invalid Page Title")
        }
        return m[2], nil 
}


//load page from the disk
func loadPage(title string) (*Page, error) {
        filename := title + ".txt"
        body, err := ioutil.ReadFile(filename)
        if err != nil {
                return  nil, err
        }
        return &Page{Title: title, Body: body}, nil
}


//helper funcion to render a template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
        err := templates.ExecuteTemplate(w, tmpl + ".html", p)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}




func handler (w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}



//view page handler 
func viewHandler (w http.ResponseWriter, r *http.Request, title string){
        p, err := loadPage(title)
        if err != nil {
                http.Redirect(w, r, "/edit/"+title, http.StatusFound)
                return 
        }
        renderTemplate(w, "view", p)
}





func editHandler(w http.ResponseWriter, r *http.Request, title string){
        p, err := loadPage(title)
        if err != nil {
                p = &Page{Title: title}
        }
        renderTemplate(w, "edit", p)
}



func saveHandler(w http.ResponseWriter, r *http.Request, title string){
        body := r.FormValue("body")
        p := &Page{Title: title, Body: []byte(body)}
        err := p.save()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return 
        }
        http.Redirect(w, r, "/view/" + title, http.StatusFound)
        
}



func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
        return func(w http.ResponseWriter, r *http.Request) {
                m := validPath.FindStringSubmatch(r.URL.Path)
                if m == nil {
                        http.NotFound(w,r)
                        return
                }
                fn(w, r, m[2])
        }

}



func main() {
       http.HandleFunc("/view/", makeHandler(viewHandler))
       http.HandleFunc("/edit/", makeHandler(editHandler))
       http.HandleFunc("/save/", makeHandler(saveHandler))
       log.Fatal(http.ListenAndServe(":8080", nil))
}
