package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //rootアクセス時chat.htmlを生成する
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(
			filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

var rooms []*room

func CreateOrJoinRoomHandler(w http.ResponseWriter, r *http.Request) {

	for roomID := 0; roomID < MAX_ROOM_NUM; roomID++ {

	}
}

func CreateRooms() {
	for i := 0; i < MAX_ROOM_NUM; i++ {
		rooms = append(rooms, newRoom())
	}
}

func CreateRoomsHTTPHandle() {
	for i := 0; i < MAX_ROOM_NUM; i++ {
		room_url := "/room" + strconv.Itoa(i)
		fmt.Println(room_url)
		http.Handle(room_url, rooms[i])
	}
}

func RunRooms() {
	for i := 0; i < MAX_ROOM_NUM; i++ {
		go rooms[i].run()
	}
}
func main() {
	//r := newRoom()
	CreateRooms()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	//http.Handle("/room", r)
	CreateRoomsHTTPHandle()
	http.HandleFunc("/createOrJonRoom", CreateOrJoinRoomHandler)
	//go r.run()
	RunRooms()
	fmt.Println("Start the ChatService")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
