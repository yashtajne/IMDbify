package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func Admin(cmd string, params string) (string, error) {
	switch cmd {
	case "deleteAll":
		_, err := GetCollection(params).DeleteMany(Ctx, bson.D{})
		if err != nil {
			return "", fmt.Errorf("error with this command: %v", err)
		}
		return "success", nil
	default:
		return "", fmt.Errorf("unknown command: %s", cmd)
	}
}
