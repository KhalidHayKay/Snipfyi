package main

import (
	"context"
	"time"
)

type Command struct {
	Run         func(ctx context.Context)
	Destructive bool
}

type Migration struct {
	Name string
	Up   string
	Down string
}

type Schema struct {
	Id        int64
	Name      string
	AppliedAt time.Time
}
