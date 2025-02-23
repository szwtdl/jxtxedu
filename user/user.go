package user

import (
	"errors"
	"github.com/szwtdl/jxtxedu/types"
	"github.com/szwtdl/jxtxedu/utils"
	"github.com/szwtdl/req"
)

func Login(client *client.HttpClient, data map[string]string) (types.User, error) {
	payload := map[string]interface{}{
		"id_card":    data["username"],
		"password":   data["password"],
		"image_code": data["image_code"],
		"sign":       "",
	}
	encryptData, _ := utils.EncryptData(payload)
	response, err := client.DoPost(types.LoginUrl, map[string]string{
		"data": encryptData,
	})
	if err != nil {
		return types.User{}, errors.New(err.Error())
	}
	var responseApi types.ResponseApi
	err = utils.JsonUnmarshal(response, &responseApi)
	if err != nil {
		return types.User{}, err
	}
	if responseApi.Code != 0 {
		return types.User{}, errors.New(responseApi.Msg)
	}
	var user types.User
	err = utils.JsonUnmarshal(utils.JsonMarshal(responseApi.Data), &user)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

func Captcha(client *client.HttpClient) (string, error) {
	response, err := client.DoGet(types.ImageCode)
	if err != nil {
		return "", errors.New(err.Error())
	}
	var responseApi types.ResponseApi
	err = utils.JsonUnmarshal(response, &responseApi)
	if err != nil {
		return "", err
	}
	if responseApi.Code != 0 {
		return "", errors.New(responseApi.Msg)
	}
	type Code struct {
		Image string `json:"image"`
	}
	var code Code
	err = utils.JsonUnmarshal(utils.JsonMarshal(responseApi.Data), &code)
	if err != nil {
		return "", err
	}
	return code.Image, nil
}
