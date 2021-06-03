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

var blogPostCollection *mgo.Collection // mongoDB collection
var commentsCollection *mgo.Collection
var repliesCollection *mgo.Collection

var tpl *template.Template

type Reply struct {
	Replier, Reply, ID, BelongsTo string
}

type Comment struct {
	Commentor, Comment, ID, BelongsTo, ReplyRoute string
	Replies                                       []Reply
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

	blogPostCollection = session.DB("golang").C("blogPost")
	commentsCollection = session.DB("golang").C("comments")
	repliesCollection = session.DB("golang").C("replies")
	//routes

	http.HandleFunc("/", Home)
	http.HandleFunc("/new-post", newPost)
	http.HandleFunc("/post/", PostSingle)
	http.HandleFunc("/reply/", replyToComment)

	log.Println("Now serving on port 8080.....")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		Check(err, "Listen and Serve")
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	var blogPosts []BlogPost

	err := blogPostCollection.Find(bson.M{}).All(&blogPosts)
	Check(err, "Finding blogposts")

	tpl.ExecuteTemplate(w, "home.html", blogPosts)
}

func newPost(w http.ResponseWriter, r *http.Request) {
	//get all Blogpost in database IDs
	var blogPosts []BlogPost
	var IDs []string

	err := blogPostCollection.Find(bson.M{}).All(&blogPosts)
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

		err = blogPostCollection.Insert(&newPost) // store newpost
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
	// fmt.Println("Post Id:",PostId)

	var bp BlogPost // get blogpost from DB using ID
	err := blogPostCollection.Find(bson.M{"id": PostId}).One(&bp)
	Check(err, "Find in PostSingle")

	var Comments []Comment
	err = commentsCollection.Find(bson.M{"belongsto": PostId}).All(&Comments)
	Check(err, "Finding comments in PostSingle")
	bp.Comments = Comments

	if r.Method == http.MethodGet {
		blogPostAndPath := postAndPath{bp, path}
		tpl.ExecuteTemplate(w, "singlePost.html", blogPostAndPath) // send data to template
	} else if r.Method == http.MethodPost { // user trying to comment
		name := r.FormValue("author")
		message := r.FormValue("message")
		rand.Seed(time.Now().UnixNano())
		ID := rand.Intn(99999) + 1
		IDAsString := strconv.Itoa(ID)

		// check if ID is not already taken
		var comment Comment
		for {
			if ID > 100000 {
				ID = 0
			}

			err := commentsCollection.Find(bson.M{"id": IDAsString}).One(&comment)
			if err != nil {
				break
			}

			ID++
		}

		replyRoute := "comment-" + IDAsString

		newcomment := Comment{Commentor: name, Comment: message, ID: IDAsString, ReplyRoute: replyRoute, BelongsTo: PostId}

		err = commentsCollection.Insert(&newcomment) // update comment in DB
		Check(err, "Updating comment")

		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}

type commentAndPath struct {
	Comment
	Path string
}

func replyToComment(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("reply path:", path)

	commentId := path[15:]
	fmt.Println("Comment id:", commentId)

	var comment Comment
	err := commentsCollection.Find(bson.M{"id": commentId}).One(&comment)
	Check(err, "Finding Comment in replytocomment")
	fmt.Println("comment:", comment)

	var RepliesToComment []Reply
	err = repliesCollection.Find(bson.M{"belongsto": commentId}).All(&RepliesToComment)
	Check(err, "Finding Replies in replytocomment")
	comment.Replies = RepliesToComment

	if r.Method == http.MethodGet {
		commentWithPath := commentAndPath{comment, path}
		tpl.ExecuteTemplate(w, "reply.html", commentWithPath)
	} else if r.Method == http.MethodPost {
		name := r.FormValue("name")
		message := r.FormValue("message")
		rand.Seed(time.Now().UnixNano())
		ID := rand.Intn(999999) + 1
		IDAsString := strconv.Itoa(ID)

		// check if ID is not already taken
		var reply Reply
		for {
			if ID > 1000000 {
				ID = 0
			}

			err := repliesCollection.Find(bson.M{"id": IDAsString}).One(&reply)
			if err != nil {
				break
			}

			ID++
		}

		newreply := Reply{Replier: name, Reply: message, ID: IDAsString, BelongsTo: commentId}

		err = repliesCollection.Insert(&newreply) // update comment in DB
		Check(err, "Updating comment")

		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}
