package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"internal/db"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type UpdateLocationData struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type GetNearbyDriversData struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var jwtSecret = []byte("supersecret")

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	userID := int(claims["user_id"].(float64))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws read error:", err)
			return
		}
		var wsMsg WSMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			continue
		}
		switch wsMsg.Event {
		case "update_location":
			var data UpdateLocationData
			marshalData(wsMsg.Data, &data)
			db.SetDriverLocation(strconv.Itoa(userID), data.Lat, data.Lng)
		case "get_nearby_drivers":
			var data GetNearbyDriversData
			marshalData(wsMsg.Data, &data)
			results, _ := db.GetNearbyDrivers(data.Lat, data.Lng, data.Radius)
			drivers := []map[string]interface{}{}
			for _, loc := range results {
				drivers = append(drivers, map[string]interface{}{
					"id":  loc.Name,
					"lat": loc.Latitude,
					"lng": loc.Longitude,
				})
			}
			resp := WSMessage{Event: "nearby_drivers", Data: drivers}
			conn.WriteJSON(resp)
		}
	}
}

func marshalData(data interface{}, out interface{}) {
	b, _ := json.Marshal(data)
	json.Unmarshal(b, out)
}
