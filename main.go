package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

const WS_ADRESS = "ws://localhost:8081" //TODO 本番かローカルかで使い分けよう
var room_num = 0

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

	if will_join_room == -1 {
		will_join_room = createRoomAndInit() // 部屋を作ってhttpHandleやgoruitineの実行などもろもろやる あと部屋番号を返す
	}

}

func createRoomAndInit() int {

	room := newRoom(room_num)
	room_url := createwsAdress(room_num)
	http.Handle(room_url, room)
	go room.run()
	go room.game.run()
	rooms = append(rooms, room)

	ret := room_num
	room_num++
	return ret
}

func createwsAdress(roomID int) string {
	return WS_ADRESS + "/room" + strconv.Itoa(roomID)
}

func main() {

	http.Handle("/", &templateHandler{filename: "loby.html"})
	http.Handle("/game", &templateHandler{filename: "game.html.tpl"})
	http.Handle("/templates/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/createOrJoinRoom", CreateOrJoinRoomHandler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
