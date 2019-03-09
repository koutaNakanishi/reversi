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

const WS_ADRESS = "ws://localhost:8081" //TODO 本番かローカルかで使い分けよう
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

func CreateOrJoinRoomHandler(w http.ResponseWriter, r *http.Request) { // /createOrJoiにアクセスした時に呼ばれる
	will_join_room := -1
	for i, room := range rooms {
		if room.GetRoomNum() == 0 {
			will_join_room = i //TODO ここは部屋の作成にするべき?
			break
		} else if room.GetRoomNum() < MAX_CONNECTION_PER_ROOM {
			will_join_room = i
			break
		}
	}

	if will_join_room == -1 {
		fmt.Println("Error:Can'tJoinProperRoom")
	}

	fmt.Fprintf(w, CreatewsAdress(will_join_room))
}

func CreatewsAdress(roomID int) string {
	return WS_ADRESS + "/room" + strconv.Itoa(roomID)
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
		go rooms[i].game.run()
	}
}
func main() {
	//r := newRoom()
	fmt.Println("Start the ChatService")
	CreateRooms()
	http.Handle("/", &templateHandler{filename: "loby.html"})
	http.Handle("/game", &templateHandler{filename: "game.html.tpl"})
	http.Handle("/templates/", http.FileServer(http.Dir("./")))
	fmt.Println("Start the ChatService")
	//http.Handle("/room", r)
	CreateRoomsHTTPHandle()
	http.HandleFunc("/createOrJoinRoom", CreateOrJoinRoomHandler)
	//go r.run()
	RunRooms()
	fmt.Println("Start the ChatService")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
