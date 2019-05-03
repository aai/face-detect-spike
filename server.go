package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/gorilla/websocket"
)

const (
	address   = ":3000"
	staticDir = "static"
)

var face *FaceDetector

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Response struct {
	FaceDetails []*rekognition.FaceDetail
	Command     string
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade HTTP connection:
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	state := FaceState{}

	for {
		mt, image, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt != websocket.BinaryMessage {
			log.Println("read expected a binary message")
			break
		}

		faces, err := face.DetectFaces(image)
		if err != nil {
			log.Println("rekognition:", err)
			break
		}

		cmd := state.AnalyzeFaces(faces.FaceDetails)

		resp := Response{
			FaceDetails: faces.FaceDetails,
			Command:     cmd,
		}

		msg, err := json.Marshal(resp)
		if err != nil {
			log.Println("marshal:", err)
			break
		}
		// log.Printf("%s\n", msg)

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	log.SetFlags(0)

	_, fakeAPI := os.LookupEnv("FAKE_API")
	face = NewFaceDetector(fakeAPI)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fs)
	http.HandleFunc("/socket", websocketHandler)

	log.Printf("Listening at %v...\n", address)
	http.ListenAndServe(address, nil)
}
