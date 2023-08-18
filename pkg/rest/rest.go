package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prondos/axdb/pkg/db"
)

type Service[IT comparable, MT any, DT any] struct {
	table *db.Table[string, MT, DT]
}

func NewService[MT any, DT any](table *db.Table[string, MT, DT]) *Service[string, MT, DT] {
	s := &Service[string, MT, DT]{
		table: table,
	}
	return s
}

func (s *Service[IT, MT, DT]) Index(c *gin.Context) {
	indices := s.table.List()
	c.IndentedJSON(http.StatusOK, indices)
}

func (s *Service[IT, MT, DT]) Get(c *gin.Context) {
	key := c.Param("key")
	rec, err := s.table.Read(key)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, rec)
}

func (s *Service[IT, MT, DT]) Put(c *gin.Context) {
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
