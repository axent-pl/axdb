package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (s *Server[IT, DT]) Start(ctx context.Context) error {
	// Initialize the Gin router.
	router := gin.Default()

	// Define routes for the REST API.
	router.GET("/items", s.service.Index)
	router.GET("/items/:key", s.service.Get)
	router.PUT("/items/:key", s.service.Put)

	// Initialize HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
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

func (s *Service[IT, DT]) Index(c *gin.Context) {
	indices := s.table.List()
	c.IndentedJSON(http.StatusOK, indices)
}

func (s *Service[IT, DT]) Get(c *gin.Context) {
	key := c.Param("key")
	rec, err := s.table.Read(key)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, rec)
}

func (s *Service[IT, DT]) Put(c *gin.Context) {
	var val DT
	key := c.Param("key")
	if err := c.BindJSON(&val); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := s.table.InsertOrUpdate(key, val); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, val)
}
