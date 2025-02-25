package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var sessions = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()
    
    var sessionID string
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            delete(sessions, sessionID)
            return
        }
        var data map[string]interface{}
        json.Unmarshal(msg, &data)
        
        if sid, ok := data["sessionId"].(string); ok {
            sessionID = sid
            sessions[sessionID] = conn
        }
        
        for id, peerConn := range sessions {
            if id != sessionID {
                peerConn.WriteMessage(websocket.TextMessage, msg)
            }
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleConnection)
    fmt.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
