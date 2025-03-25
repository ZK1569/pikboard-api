package client

import "sync"

var lock *sync.Mutex

func init() {
	lock = &sync.Mutex{}
}

type ChessImage interface {
	FemToImage(string, bool) ([]byte, error)
}
