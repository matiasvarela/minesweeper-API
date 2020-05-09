package random

import (
	"math/rand"
	"time"
)

type Random interface {
	Init()
	GenerateN(n int) []int
}

type random struct{}

func NewRandom() Random {
	return &random{}
}

func (r *random) Init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (r *random) GenerateN(n int) []int {
	return rand.Perm(n)
}
