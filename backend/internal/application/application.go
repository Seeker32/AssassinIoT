package application

import (
	"log/slog"

	"github.com/Seeker32/AssassinIoT/backend/internal/dependence"
)

type server struct {
	dep    dependence.Dep
	logger *slog.Logger
}

func NewServer(dep dependence.Dep) *server {
	return &server{
		dep:    dep,
		logger: dep.Logger(),
	}
}
