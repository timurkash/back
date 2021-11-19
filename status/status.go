package status

import (
	"errors"
	"fmt"
	"gitlab.com/mcsolutions/lib/back/common/json"
	"gitlab.com/mcsolutions/lib/back/common/logger"
	"log"
	"net/http"
)

type (
	Status struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
	}
	StatusOut struct {
		Status *Status `json:"status"`
	}
	IdOut struct {
		Id int64 `json:"id"`
	}
)

var (
	StatusOk              = Status{Ok: true}
	statusOk              = []byte(`{"status": {"ok": true}}`)
	statusJsonEncodeError = []byte(`{"status": {"ok": false, "error": "json encode error"}}`)
)

func GetStatusOut(err error) *StatusOut {
	s := err == nil
	statusOut := StatusOut{}
	if s {
		statusOut = StatusOut{&Status{Ok: s}}
	} else {
		log.Println(err)
		statusOut = StatusOut{&Status{Ok: s, Error: err.Error()}}
	}
	return &statusOut
}

func GetStatusOutLogger(err error, logger *logger.Logger) *StatusOut {
	s := err == nil
	statusOut := StatusOut{}
	if s {
		statusOut = StatusOut{&Status{Ok: s}}
	} else {
		log.Println(err)
		if logger != nil {
			logger.Zap.Error(err.Error())
			logger.SendError(err)
		}
		statusOut = StatusOut{&Status{Ok: s, Error: err.Error()}}
	}
	return &statusOut
}

func WriteStatus(w *http.ResponseWriter, err error) {
	if err := json.Encode(GetStatusOut(err), *w); err != nil {
		StatusJsonEncodeError(w, err, nil)
	}
}

func WriteLoggerStatus(w *http.ResponseWriter, err error, logger *logger.Logger) {
	if err := json.Encode(GetStatusOut(err), *w); err != nil {
		StatusJsonEncodeError(w, err, logger)
	}
}

func ErrorStatusOut(err error) []byte {
	bytes, errM := json.EncodeBytes(
		&StatusOut{
			&Status{
				Ok:    false,
				Error: err.Error(),
			},
		},
	)
	if errM != nil {
		return statusJsonEncodeError
	}
	return bytes
}

func StatusJsonEncodeError(w *http.ResponseWriter, err error, logger *logger.Logger) {
	(*w).WriteHeader(http.StatusUnprocessableEntity)
	(*w).Write(statusJsonEncodeError)
	if logger != nil {
		logger.Zap.Error(err.Error())
		logger.SendError(err)
	}
}

func StatusOutBytes(inter interface{}, err error) ([]byte, error) {
	if err != nil {
		//bytes, errM := json.EncodeByte(&StatusOut{&Status{Ok: false, Error: err.Error()}})
		//if errM != nil {
		//	return &statusJsonEncodeError, errM
		//}
		return ErrorStatusOut(err), err
	}
	if inter == nil {
		return statusOk, nil
	}

	bytes, errM := json.EncodeBytes(inter)
	if errM != nil {
		return statusJsonEncodeError, errM
	}
	//bytes[len(bytes)-1] = ','
	//statusBytes, errM := json.Marshal(statusOut)
	//if errM != nil {
	//	return &statusInternalError, errM
	//}
	//bytes = append(bytes, statusBytes[1:len(statusBytes)]...)
	return bytes, nil
}

func StatusOutCodeBytes(inter interface{}, err error) (int, []byte, error) {
	if err != nil {
		//bytes, errM := json.EncodeByte(&StatusOut{&Status{Ok: false, Error: err.Error()}})
		//if errM != nil {
		//	return http.StatusUnprocessableEntity, &statusJsonEncodeError, errM
		//}
		return http.StatusBadRequest, ErrorStatusOut(err), err
	}
	if inter == nil {
		return http.StatusOK, statusOk, nil
	}
	bytes, errM := json.EncodeBytes(inter)
	if errM != nil {
		return http.StatusUnprocessableEntity, statusJsonEncodeError, errM
	}
	return http.StatusOK, bytes, nil
}

func GetStatus(resp *http.Response) error {
	if !CodeOK(resp) {
		return fmt.Errorf("status code %d is not OK", resp.StatusCode)
	}
	status_ := &StatusOut{}
	if err := json.Decode(resp.Body, status_); err != nil {
		return err
	}
	if status_.Status == nil {
		return nil
	}
	if !status_.Status.Ok {
		return errors.New(status_.Status.Error)
	}
	return nil
}

func CodeOK(resp *http.Response) bool {
	code := resp.StatusCode
	return code == http.StatusOK ||
		code == http.StatusCreated ||
		code == http.StatusAccepted ||
		code == http.StatusNonAuthoritativeInfo
}
