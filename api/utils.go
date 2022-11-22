package api

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenRandString(length int) string {
	r := make([]rune, length)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	return string(r)
}
