package api

import "github.com/christopherplain/chirpy/internal/model"

type ApiConfig struct {
	DB             *model.DB
	FileserverHits int
	JwtSecret      string
	PolkaKey       string
}
