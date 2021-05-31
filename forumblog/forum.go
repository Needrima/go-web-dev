package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	// "strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var forumCollection *mgo.Collection // mongoDB collection

var tpl *template.Template

type Reply struct {
	Author, Reply string
}

type Comment struct {
	Author, Comment, ID, ReplyRoute string
	Replies                         []string
}

type BlogPost struct {
	Author, Title, Message, ID, SingleRoute string
	Comments                                []Comment
}

func Check(err error, msg string) {
	if err != nil {
		fmt.Println(msg, ":", err)
		return
	}
}

func Found(data []string, datum string) bool {
	for _, v := range data {
		if datum == v {
			return true
		}
	}
	return false
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	Check(err, "Session")
	defer session.Close()

	err = session.Ping()
	Check(err, "Connection to database failed")

	forumCollection = session.DB("golang").C("go_practice")

	//routes

	http.HandleFunc("/", Home)
	http.HandleFunc("/new-post", newPost)
	http.HandleFunc("/post/", PostSingle)
	//http.HandleFunc("/reply/", Reply)

	log.Println("Now serving on port 8080.....")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		Check(err, "Listen and Serve")
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	var blogPosts []BlogPost

	err := forumCollection.Find(bson.M{}).All(&blogPosts)
	Check(err, "Finding blogposts")

	tpl.ExecuteTemplate(w, "home.html", blogPosts)
}

func newPost(w http.ResponseWriter, r *http.Request) {
	//get all Blogpost in database IDs
	var blogPosts []BlogPost
	var IDs []string

	err := forumCollection.Find(bson.M{}).All(&blogPosts)
	Check(err, "Find in newpost")

	for _, v := range blogPosts {
		IDs = append(IDs, v.ID)
	}

	if r.Method == http.MethodGet { // execute template
		tpl.ExecuteTemplate(w, "post.html", "Create thread")
	} else if r.Method == http.MethodPost {
		author := r.FormValue("author")
		title := r.FormValue("title")
		message := r.FormValue("message")
		rand.Seed(time.Now().UnixNano())
		ID := rand.Intn(999) + 1       // ID for new post
		IDAsString := strconv.Itoa(ID) // convert to string

		for { // check if ID is already taken
			if !Found(IDs, IDAsString) {
				break
			} else {
				if ID > 1000 {
					ID = 0
				}
			}
			ID++
		}

		singleRoute := "post-" + IDAsString // route to view post single

		newPost := BlogPost{Author: author, Title: title, Message: message, ID: IDAsString, SingleRoute: singleRoute} // create new post

		err = forumCollection.Insert(&newPost) // store newpost
		Check(err, "Inserting blogpost")

		tpl.ExecuteTemplate(w, "post.html", "Thread Created")
	}
}

type postAndPath struct {
	BlogPost
	Path string
}

func PostSingle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path // get path

	PostId := path[11:] // trim to get post Id
	fmt.Println(PostId)

	var bp BlogPost // get blogpost from DB using ID
	err := forumCollection.Find(bson.M{"id": PostId}).One(&bp)
	Check(err, "Find in PostSingle")

	if r.Method == http.MethodGet {
		blogPostAndPath := postAndPath{bp, path}
		tpl.ExecuteTemplate(w, "singlePost.html", blogPostAndPath) // send data to template
	} else if r.Method == http.MethodPost { // user trying to comment
		name := r.FormValue("author")
		message := r.FormValue("message")

		// get former comments ID
		var blogPosts []BlogPost
		var CommentIDs []string

		err := forumCollection.Find(bson.M{}).All(&blogPosts)
		Check(err, "Find in postSingle")

		for _, v := range blogPosts {
			for _, v2 := range v.Comments {
				CommentIDs = append(CommentIDs, v2.ID)
			}
		}
		//define ID for new comment
		rand.Seed(time.Now().UnixNano())
		ID := rand.Intn(99999) + 1
		IDAsString := strconv.Itoa(ID)

		// check if ID is not already taken
		for {
			if !Found(CommentIDs, IDAsString) {
				break
			} else {
				if ID > 100000 {
					ID = 0
				}
			}
			ID++
		}

		replyRoute := "comment-" + IDAsString

		comment := Comment{Author: name, Comment: message, ReplyRoute: replyRoute}

		bp.Comments = append(bp.Comments, comment) // update comment

		err = forumCollection.Update(bson.M{"id": PostId}, bson.M{"$set": bson.M{"comments": bp.Comments}}) // update comment in DB
		Check(err, "Updating comment")

		blogPostAndPath := postAndPath{bp, path}
		tpl.ExecuteTemplate(w, "singlePost.html", blogPostAndPath)
	}
}
