package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_038_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	//log := logger.GetLogger(r.Context())
	//log.WithFields(logrus.Fields{"body": string(body)}).Debugf("request body")

	err = json.Unmarshal(body, dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}
