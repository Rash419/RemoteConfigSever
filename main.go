package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
  log.Println("Starting RemoteConfigServer...")

	http.HandleFunc("/template.json", templateConfigHandler)
	http.Handle("/static/", http.StripPrefix("/static/impress-template", http.FileServer(http.Dir("./web/static/impress_template/"))))

  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Print("Failed to start RemoteConfigServer")
  }
}

type uriJSONObj struct {
	Uri   string `json:"uri"`
	Stamp string `json:"stamp"`
}

type remoteTemplateJSON struct {
	Kind      string     `json:"kind"`
	Server    string     `json:"server"`
	Templates []uriJSONObj `json:"templates"`
}

func getTemplateConfigJSON(kind string, server string, templates []uriJSONObj) ([]byte, error) {
	data := &remoteTemplateJSON{
		Kind:      kind,
		Server:    server,
		Templates: templates,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func templateConfigHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Recieved template.json request");
  const uriFormat = "http://localhost:8080/static/impress-template/template%d.otp";
  const kind = "templateconfiguration"
  const server = "remoteserver"
  uriObjArr := make([]uriJSONObj, 3)
  for i := 0; i < 3; i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(uriFormat, i+1),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    uriObjArr[i] = *uriObj
  }
  templateJSON, err := getTemplateConfigJSON("templateconfiguration", server, uriObjArr)
  if err != nil {
    log.Printf("Failed to parse template config JSON with err[%s]", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(templateJSON)
  if err != nil {
    log.Printf("Failed to write json with error[%s]", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}
