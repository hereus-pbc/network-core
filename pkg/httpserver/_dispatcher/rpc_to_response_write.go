package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func rpcToResponseWriteJSON(output interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func rpcToResponseWrite(output interface{}, w http.ResponseWriter) {
	switch outputType := output.(type) {
	case map[string]interface{}:
		rpcToResponseWriteJSON(outputType, w)
	case []interface{}:
		rpcToResponseWriteJSON(outputType, w)
	case struct{}:
		rpcToResponseWriteJSON(outputType, w)
	case nil:
		w.WriteHeader(http.StatusOK)
		return
	case string:
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(outputType)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	case []byte:
		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err := w.Write(outputType); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	case int:
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(fmt.Sprintf("%d", outputType))); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	case float64:
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(fmt.Sprintf("%f", outputType))); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	case bool:
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(fmt.Sprintf("%t", outputType))); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}

func convertRpcResponseToHttpResponse(output interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rpcToResponseWrite(output, w)
}
