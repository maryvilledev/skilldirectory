package controller

import (
	"net/http"
	"skilldirectory/data"

	"github.com/Sirupsen/logrus"
)

type BaseController struct {
	w http.ResponseWriter
	r *http.Request
	*logrus.Logger
	session    data.DataAccess
	fileSystem data.FileSystem
}

func (bc *BaseController) Init(w http.ResponseWriter, r *http.Request,
	session data.DataAccess, fs data.FileSystem, logger *logrus.Logger) {
	bc.w = w
	bc.r = r
	bc.Logger = logger
	bc.session = session
	bc.fileSystem = fs
}
