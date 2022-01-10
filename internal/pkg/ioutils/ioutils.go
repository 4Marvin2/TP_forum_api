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
	// if err != nil {
	// 	SendError(w, http.StatusInternalServerError, models.ModelError{
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }
}

func SendError(w http.ResponseWriter, respCode int, errorMsg string) {
	// w.WriteHeader(respCode)
	// _ = WriteJSON(w, errorMsg)
	Send(w, respCode, models.ModelError{
		Message: errorMsg,
	})
}

func SendWithoutBody(w http.ResponseWriter, respCode int) {
	w.WriteHeader(respCode)
	// if err != nil {
	// 	SendError(w, http.StatusInternalServerError, models.ModelError{
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }
}

// func ReadJSON(r *http.Request, data interface{}) error {
// 	byteReq, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal(byteReq, &data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func WriteJSON(w http.ResponseWriter, data interface{}) error {
// 	byteResp, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Write(byteResp)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
