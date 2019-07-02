package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Contains(arr primitive.A, x string) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}
