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
var room_num = 0

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

func CreateOrJoinRoomHandler(w http.ResponseWriter, r *http.Request) { // /createOrJoinにアクセスした時に呼ばれる
	will_join_room := -1
	for _, room := range rooms {
		if room.GetRoomNum() == 0 {
			will_join_room = room.room_id //TODO ここは部屋の作成にするべき?
			break
		} else if room.GetRoomNum() < MAX_CONNECTION_PER_ROOM {
			will_join_room = room.room_id
			break
		}
	}
	//fmt.Println("willjoinROom", will_join_room)
	if will_join_room == -1 {
		//rooms = append(rooms, createRoom(len(rooms)))
		will_join_room = createRoomAndInit() // 部屋を作ってhttpHandleやgoruitineの実行などもろもろやる あと部屋番号を返す
		//fmt.Println("room拡張!今roomnum:", room_num)
	}

	fmt.Fprintf(w, CreatewsAdress(will_join_room))
}

func createRoomAndInit() int {

	room := newRoom(room_num)
	room_url := "/room" + strconv.Itoa(room_num)
	http.Handle(room_url, room)
	go room.run()
	go room.game.run()
	rooms = append(rooms, room)

	ret := room_num
	room_num++
	return ret
}

func CreatewsAdress(roomID int) string {
	return WS_ADRESS + "/room" + strconv.Itoa(roomID)
}

/*
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
*/
func test() {
	for {
		var a string
		fmt.Scan(&a)
		if a == "rooms" {
			for _, i := range rooms {
				fmt.Println(i.room_id, len(i.clients))
			}
		}
	}
}
func main() {
	//r := newRoom()
	go test()
	//CreateRooms()
	http.Handle("/", &templateHandler{filename: "loby.html"})
	http.Handle("/game", &templateHandler{filename: "game.html.tpl"})
	http.Handle("/templates/", http.FileServer(http.Dir("./")))
	//http.Handle("/room", r)
	//CreateRoomsHTTPHandle()
	http.HandleFunc("/createOrJoinRoom", CreateOrJoinRoomHandler)
	//go r.run()
	//RunRooms()
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
