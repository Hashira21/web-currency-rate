package tech

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var (
	appInfo *app
)

func New() *tech {
	return &tech{
		app: &app{},
	}
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	writeJson(w, appInfo)
}

func GetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy"}`))
}

func (t *tech) SetAppInfo(name, version string) *tech {
	t.app.Name, t.app.Version = name, version
	t.getGuidAndHostname()
	appInfo = t.app

	return t
}

func (t *tech) getGuidAndHostname() {
	t.app.Guid = uuid.New().String()

	host, err := os.Hostname()
	if err != nil {
		t.logger.Warn().Msg(fmt.Sprintf("failed to get hostname for %s", t.app.Name))
	} else {
		t.app.Hostname = host
	}
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if dataBytes, err := json.Marshal(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(dataBytes)
	}
}
