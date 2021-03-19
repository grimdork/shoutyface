package main

import (
	"context"
	"time"

	"github.com/Urethramancer/signor/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
	"github.com/grimdork/sweb"
	"github.com/jackc/pgx/v4/pgxpool"
)

// API endpoints:
// POST /message - token, sender, destinations, subject, template name, JSON data

// Shoutyface runtime data.
type Shoutyface struct {
	sweb.Server
	dbp *pgxpool.Pool

	mailquit  chan interface{}
	mailqueue chan Message
	mailhost  string
	mailport  string
	mailuser  string
	mailpass  string
	mailfrom  string
}

// NewServer initialises the API server.
func NewServer() (*Shoutyface, error) {
	srv := &Shoutyface{
		mailquit:  make(chan interface{}),
		mailqueue: make(chan Message),
		mailhost:  env.Get("MAILHOST", "localhost"),
		mailport:  env.Get("MAILPORT", "587"),
		mailuser:  env.Get("MAILUSER", ""),
		mailpass:  env.Get("MAILPASS", ""),
		mailfrom:  env.Get("MAILFROM", "root@localhost"),
	}

	srv.Init()
	var err error
	srv.dbp, err = pgxpool.Connect(context.Background(), env.Get("DATABASE_URL", "localhost"))
	if err != nil {
		srv.E("Database error: %s", err.Error())
		return nil, err
	}

	var id int
	row := srv.dbp.QueryRow(context.Background(), "select id from users limit 1")
	err = row.Scan(&id)
	if err != nil && err.Error() != "" {
		err := srv.createTables()
		if err != nil {
			return nil, err
		}
	}

	srv.L("Connected to database.")
	srv.InitMiddleware()
	srv.AddStartHook(srv.StartShouty)
	srv.AddStopHook(srv.StopShouty)

	srv.Route("/api/message", func(r chi.Router) {
		r.Use(httprate.LimitByIP(60, 1*time.Minute))
		r.Use(sweb.AddJSONHeaders)
		r.Use(srv.tokencheck)
		r.Post("/", srv.PostMessage)
	})

	srv.Route("/api", func(r chi.Router) {
		r.Use(httprate.LimitByIP(60, 1*time.Minute))
		r.Use(sweb.AddJSONHeaders)
		r.Use(srv.admintokencheck)

		r.Get("/users", srv.GetUsers)
		r.Post("/user", srv.PostUser)
		r.Delete("/user", srv.DeleteUser)

		r.Post("/subscribe", srv.PostSubscribe)
		r.Delete("/subscribe", srv.DeleteSubscription)
		r.Get("/subscriptions", srv.GetSubscriptions)

		r.Post("/channel", srv.PostChannel)
		r.Delete("/channel", srv.DeleteChannel)
		r.Get("/channels", srv.GetChannels)
	})
	return srv, nil
}

// StartShouty server.
func (sh *Shoutyface) StartShouty() error {
	sh.L("Servcer starting.")
	go func() {
		sh.RunMailQueue(sh.mailquit)
	}()
	return nil
}

// StopShouty gracefully.
func (srv *Shoutyface) StopShouty() {
	srv.mailquit <- true
	srv.dbp.Close()
	srv.Wait()
	srv.L("Server stopped")
}
