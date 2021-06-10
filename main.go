package main

import (
	"context"

	"notes-api/api"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	if err := api.ListenAndServe(ctx); err != nil {
		logrus.WithContext(ctx).WithError(err).Error("Error starting API server")
	}
}
