package controller

import (
	"context"
	"net/http"
)

type Cron struct {
}

func NewCronController() *Cron {
	return &Cron{}
}

// Show endpoint
func (c *Cron) Show(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	return http.StatusOK, "a", nil
}
