package models

import (
	"time"
)

type Resources struct {
	Id         int64     `json:"id" bun:"id,pk,autoincrement"`
	ResourceId string    `json:"resourceId" bun:"resource_id,unique"`
	Account    string    `json:"account"`
	AccountId  string    `json:"accountId" bun:"account_id"`
	Service    string    `json:"service"`
	Region     string    `json:"region"`
	Name       string    `json:"name"`
	Cost       float64   `json:"cost"`
	CreatedAt  time.Time `json:"createdAt" bun:"created_at"`
	FetchedAt  time.Time `json:"fetchedAt" bun:"fetched_at"`
}

type Edges struct {
	Source int64  `json:"source" bun:"source"`
	Dest   int64  `json:"dest" bun:"dest"`
	Name   string `json:"name" bun:"name"`
}
