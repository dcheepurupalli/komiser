package models

import (
	"time"
)

type Resource struct {
	Id         int64             `json:"id" bun:"id,pk,autoincrement"`
	ResourceId string            `json:"resourceId" bun:"resource_id,unique"`
	Provider   string            `json:"provider"`
	Account    string            `json:"account"`
	AccountId  string            `json:"accountId" bun:"account_id"`
	Service    string            `json:"service"`
	Region     string            `json:"region"`
	Name       string            `json:"name"`
	CreatedAt  time.Time         `json:"createdAt" bun:"created_at"`
	FetchedAt  time.Time         `json:"fetchedAt" bun:"fetched_at"`
	Cost       float64           `json:"cost"`
	Metadata   map[string]string `json:"metadata"`
	Tags       []Tag             `json:"tags" bun:"tags,default:'[]'"`
	Link       string            `json:"link" bson:"link"`
	Data       string            `json:"data" bun:"data,default:'{}'"`
	Value      string            `bun:",scanonly"` // to be deprecated
	Secrets    []KeyValuePair    `json:"secrets" bun:"secrets,default:'[]'"`
	Variables  []KeyValuePair    `json:"variables" bun:"variables,default:'[]'"`
	SBOM       string            `json:"sbom" bun:"sbom,default:'{}'"`
	Unmanaged  string            `json:"gitleaks" bun:"gitleaks,default:'[]'"`
}

type Tag struct {
	Key   string `json:"key" bun:"key"`
	Value string `json:"value" bun:"value"`
}

type BulkUpdateTag struct {
	Tags      []Tag `json:"tags"`
	Resources []int `json:"resources"`
}

type KeyValuePair struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
}
