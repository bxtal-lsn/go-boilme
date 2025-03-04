package boilme

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func (b *Boilme) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     b.ErrorLog,
		Handler:      b.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if b.DB.Pool != nil {
		defer b.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerConn != nil {
		defer badgerConn.Close()
	}

	go b.listenRPC()
	b.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	return srv.ListenAndServe()
}
