package main

import (
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
)

type reply struct {
	Chat, Reply string
}

var chatdb = map[string]string{}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.html", "reply.html"))
}

func Chat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Chat")
	if err != nil {
		chatID, _ := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "Chat",
			Value:    chatID.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	var chtext string
	if r.Method == http.MethodPost {
		chtext = r.FormValue("chat")
		chatdb[cookie.Value] = chtext
	}
	tpl.ExecuteTemplate(w, "index.html", chtext)
}

func Reply(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Chat")
	if err != nil {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	cht, ok := chatdb[cookie.Value]
	if !ok {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "reply.html", reply{cht, ""})
	} else if r.Method == http.MethodPost {
		rply := r.FormValue("reply")
		r := reply{cht, rply}
		tpl.ExecuteTemplate(w, "reply.html", r)
	}

}

func main() {
	http.HandleFunc("/chat", Chat)
	http.HandleFunc("/reply", Reply)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
