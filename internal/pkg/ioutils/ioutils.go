package ioutils

import (
	"forumApp/internal/forumapp/models"
	"io/ioutil"
	"net/http"
)

type ReadModel interface {
	UnmarshalJSON(data []byte) error
}

type WriteModel interface {
	MarshalJSON() ([]byte, error)
}

func Send(w http.ResponseWriter, respCode int, respBody WriteModel) {
	w.WriteHeader(respCode)
	_ = WriteJSON(w, respBody)
}

func SendError(w http.ResponseWriter, respCode int, errorMsg string) {
	Send(w, respCode, models.ModelError{
		Message: errorMsg,
	})
}

func SendWithoutBody(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
}

func ReadJSON(r *http.Request, data ReadModel) error {
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = data.UnmarshalJSON(byteReq)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, data WriteModel) error {
	byteResp, err := data.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = w.Write(byteResp)
	if err != nil {
		return err
	}

	return nil
}
