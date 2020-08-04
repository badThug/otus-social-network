package web

import (
	"github.com/badThug/otus-social-network/app/globals"
	"github.com/badThug/otus-social-network/app/handlers"
)

func initRoutes(d Dispatcher) {

	// IndexHandler
	d.Get(globals.ViewIndexRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ViewIndexHandler()
	}))

	// FriendsHandler
	d.Post(globals.ApiFriendsRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiAddFriendHandler()
	}))
	d.Get(globals.ApiFriendsRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiGetFriendsHandler()
	}))
	d.Post(globals.ApiRemoveFriendsRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiDeleteFriendHandler()
	}))
	d.Get(globals.ApiAvailableFriendRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiGetAvailableFriendsHandler()
	}))

	// AuthHandler
	d.Post(globals.ApiLoginRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiLoginHandler()
	}))
	d.Post(globals.ApiRegisterRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiRegisterHandler()
	}))
	d.Post(globals.ApiLogoutRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiLogoutHandler()
	}))
	d.Get(globals.ApiGetUserRoute, d.handleRequest(func(h *handlers.Handler) error {
		return h.ApiGetUserHandler()
	}))

	d.Router.Use(SessionAuthentication)
}
