package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/frosttlink/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductService
	Sessions       *scs.SessionManager
	WsUpgradedr    websocket.Upgrader
}
