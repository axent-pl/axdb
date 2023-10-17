package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/prondos/axdb/pkg/db"
)

type Service[IT comparable, DT any] struct {
	table *db.Table[string, DT]
}

type Server[IT comparable, DT any] struct {
	service *Service[IT, DT]
}

func NewServer[DT any](service *Service[string, DT]) *Server[string, DT] {
	server := &Server[string, DT]{service: service}
	return server
}

func (s *Server[IT, DT]) Start(ctx context.Context, serverAddress string) error {
	router := &Router{}
	router.GET("^/items$", s.service.Index)
	router.GET("^/items/[^/]+$", s.service.Get)
	router.PUT("^/items/[^/]+$", s.service.Put)

	// Initialize HTTP server
	httpServer := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	// Initialize channels handling HTTP server shutdown
	done := make(chan error)

	// Start HTTP server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			done <- err
		}
		done <- nil
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("HTTP server stopped by context\n")
			return httpServer.Shutdown(context.Background())
		case err := <-done:
			log.Printf("HTTP server stopped by itself: %v\n", err)
			return err
		}
	}
}

func NewService[DT any](table *db.Table[string, DT]) *Service[string, DT] {
	s := &Service[string, DT]{
		table: table,
	}
	return s
}

func (s *Service[IT, DT]) getKey(r *http.Request) string {
	pathParts := strings.Split(r.URL.Path, "/")
	return pathParts[len(pathParts)-1]
}

func (s *Service[IT, DT]) Index(w http.ResponseWriter, r *http.Request) {
	indices := s.table.List()
	indicesJson, err := json.Marshal(indices)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(indicesJson)
}

func (s *Service[IT, DT]) Get(w http.ResponseWriter, r *http.Request) {
	key := s.getKey(r)
	data, err := s.table.Read(key)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(dataJson)
}

func (s *Service[IT, DT]) Put(w http.ResponseWriter, r *http.Request) {
	var data *DT = new(DT)
	key := s.getKey(r)
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := s.table.InsertOrUpdate(key, *data); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	s.Get(w, r)
}
