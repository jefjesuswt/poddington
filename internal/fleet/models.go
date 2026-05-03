package fleet

import (
	"errors"
	"time"
)

var (
	ErrNodeNotFound             = errors.New("node not found")
	ErrNodeNameAlreadyExists    = errors.New("node name already exists")
	ErrNodeAddressAlreadyExists = errors.New("node address already exists")
)

type Node struct {
	ID        string
	Name      string
	Address   string
	Token     string
	CreatedAt time.Time
	LastSeen  time.Time
}

type NodeTelemtry struct {
	CpuPercentUsage float64 `json:"cpu_percent_usage"`
	MemTotalMB      uint64  `json:"mem_total_mb"`
	MemUsedMB       uint64  `json:"mem_used_mb"`
	MemPercentUsage float64 `json:"mem_percent_usage"`
}
