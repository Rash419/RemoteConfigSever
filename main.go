package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
  log.Println("Starting RemoteConfigServer...")

	http.HandleFunc("/asset.json", templateConfigHandler)
  http.HandleFunc("/font.json", fontConfigHandler)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))

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
	Templates map[string][]uriJSONObj `json:"templates"`
  Fonts     []uriJSONObj `json:"fonts"`
}

func getAssetConfigJSON(kind string, server string, templates map[string][]uriJSONObj, fonts []uriJSONObj) ([]byte, error) {
	data := &remoteTemplateJSON{
		Kind:      kind,
		Server:    server,
		Templates: templates,
    Fonts: fonts,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func getFontConfigJSON(kind string, server string, fonts []uriJSONObj) ([]byte, error) {
	data := &remoteTemplateJSON{
		Kind:      kind,
		Server:    server,
    Fonts: fonts,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}


func templateConfigHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Recieved asset.json request");
  const templateUriFormat = "http://localhost:8080/static/impress-template/template%d.otp";
  const fontUriFormat = "http://localhost:8080/static/font/font%d.ttf";
  const kind = "assetconfiguration"
  const server = "remoteserver"

  templateMap := make(map[string][]uriJSONObj);

  presntTemplates := make([]uriJSONObj, 3)
  for i := 0; i < 3; i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(templateUriFormat, i+1),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    presntTemplates[i] = *uriObj
  }

  fontArr := make([]uriJSONObj, 3)
  for i := 0; i < 3; i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(fontUriFormat, i+1),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    fontArr[i] = *uriObj
  }

  templateMap["presnt"] = presntTemplates;

  templateJSON, err := getAssetConfigJSON(kind, server, templateMap, fontArr)
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

func fontConfigHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("Recieved font.json request");
  const fontUriFormat = "http://localhost:8080/static/font/font%d.ttf";
  const kind = "fontconfiguration"
  const server = "remoteserver"

  fontArr := make([]uriJSONObj, 3)
  for i := 0; i < 3; i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(fontUriFormat, i+1),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    fontArr[i] = *uriObj
  }

  fontJSON, err := getFontConfigJSON(kind, server, fontArr)
  if err != nil {
    log.Printf("Failed to parse font config JSON with err[%s]", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(fontJSON)
  if err != nil {
    log.Printf("Failed to write json with error[%s]", err.Error())
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}
