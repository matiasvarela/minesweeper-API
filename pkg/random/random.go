package random

import (
	"math/rand"
	"time"
	"github.com/google/uuid"
)

type Random interface {
	Init()
	GenerateN(n int) []int
	GenerateID() string
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

func (r *random) GenerateID() string {
	return uuid.New().String()
}
