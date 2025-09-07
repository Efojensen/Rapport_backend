package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Efojensen/rapport.git/constants"
	"github.com/gofiber/fiber/v2"
)

func JoinCommunity(userId string) (string, error) {
	communityUrl := fmt.Sprint(constants.PubUrl, userId)

	proxy := fiber.AcquireAgent()
	proxy.Request().Header.SetMethod("POST")
	proxy.Request().SetRequestURI(communityUrl)

	err := proxy.Parse()
	if err != nil {
		return "", err
	}

	var resBody struct {
		Message string `json:"message"`
	}

	_, body, errs := proxy.Bytes()
	if len(errs) > 0 {
		return "", err
	}

	err = json.Unmarshal(body, &resBody)
	if err != nil {
		return "", err
	}

	return resBody.Message, nil
}