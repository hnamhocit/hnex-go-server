package handlers

import (
	"hnex_server/internal/repositories"
)

type ProfileHandler struct {
	Repo     *repositories.ProfileRepository
	AuthRepo *repositories.AuthRepository
}
