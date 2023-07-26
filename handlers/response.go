package middlewares

import (
	"encoding/json"
	"net/http"

	"blockparty.co/test/models"
	"github.com/fatih/color"
)

func SuccessArrRespond(fields interface{}, modelType string, writer http.ResponseWriter) {
	_, err := json.Marshal(fields)
	type data struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	temp := &data{Data: fields, Message: "success"}
	if err != nil {
		ServerErrResponse(err.Error(), writer)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	switch modelType {
	case "Metadata":
		temp.Data = fields.([]*models.Metadata)
	default:
		// handle invalid model type
	}

	json.NewEncoder(writer).Encode(temp)
}

// SuccessMessageResponse -> success error messageformatter
func SuccessMessageResponse(msg string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	temp := &errdata{Message: "success", Status: msg}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(temp)
}

// ErrorResponse -> error formatter
func ErrorResponse(error string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
	}
	temp := &errdata{Message: error}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(temp)
}

func ServerErrResponse(error string, writer http.ResponseWriter) {
	type errdata struct {
		Message string `json:"message"`
	}
	temp := &errdata{Message: error}

	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(temp)
}

func SuccessOneRespond(fields interface{}, modelType string, writer http.ResponseWriter) {
	_, err := json.Marshal(fields)
	type data struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	temp := &data{Data: fields, Message: "success"}
	if err != nil {
		color.Red("Marshal Data Failed in SuccessOneRespond() for Type(%v)...", modelType)
		ServerErrResponse(err.Error(), writer)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	switch modelType {
	case "Metadata":
		temp.Data = fields.(models.Metadata)
	default:
		color.Red("Invalid Model Type in SuccessOneRespond() for Type ( %v )...", modelType)
	}

	json.NewEncoder(writer).Encode(temp)
}
