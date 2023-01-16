package response

import (
	"context"
	"fmt"
	"net/http"
)

func Error(w http.ResponseWriter, status int, message interface{}) {
	err := JSON(w, status, Envelope{"success": false, "error": message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	//logger.Errorf(ctx, "Server error - method: %s, url: %s, error: %s", r.Method, r.URL.String(), err)
	Error(w, http.StatusInternalServerError, "Internal server error")
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusNotFound, "The request resource was not found")
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this endpoint", r.Method)
	Error(w, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	err := JSON(w, http.StatusBadRequest, Envelope{"fieldErrors": errs})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusUnauthorized, "invalid authentication credentials")
}

func InvalidAuthenticationtokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	Error(w, http.StatusUnauthorized, "invalid or missing authentication token")
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

func NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusForbidden, "Your account does not have the necessary permissions to access this resource")
}

func InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	Error(w, http.StatusForbidden, "Your account must be activated to access this resource")
}
