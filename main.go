package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	Stamp string `json:"version"`
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

func getStaticFileName(kind string) ([]string, error) {
  pathToStatic := "./web/static/";
  filePaths := make([]string,0)

  switch kind {
  case "impress-template":
    pathToStatic += "impress-template"
  case "font":
    pathToStatic += "font"
  default:
    return filePaths, nil;
  }

  err := filepath.Walk(pathToStatic, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err;
    }

    if !info.IsDir() {
        filePaths = append(filePaths, info.Name())
    }
    return nil
  })

  return filePaths, err;
}

func templateConfigHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved asset.json request")
	const uriFormat = "http://localhost:8080/static/%s/%s"
	const kind = "assetconfiguration"
	const server = "remoteserver"

	templateFileNames, err := getStaticFileName("impress-template")
	if err != nil {
		log.Printf("Failed to get template static files with error[%s]", err.Error())
	}

	templateMap := make(map[string][]uriJSONObj)

	presntTemplates := make([]uriJSONObj, 0)
	for i := 0; i < len(templateFileNames); i++ {
		uriObj := &uriJSONObj{
			Uri:   fmt.Sprintf(uriFormat, "impress-template", templateFileNames[i]),
			Stamp: fmt.Sprintf("%d", i+1),
		}
		presntTemplates = append(presntTemplates, *uriObj)
	}

  fontFileNames , err := getStaticFileName("font")
  if err != nil {
		log.Printf("Failed to get font static files with error[%s]", err.Error())
  }

  fontArr := make([]uriJSONObj, 0)
  for i := 0; i < len(fontFileNames); i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(uriFormat, "font", fontFileNames[i]),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    fontArr = append(fontArr, *uriObj);
  }

  templateMap["presentation"] = presntTemplates;

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

	const uriFormat = "http://localhost:8080/static/%s/%s"
  const kind = "fontconfiguration"
  const server = "remoteserver"

  fontFileNames , err := getStaticFileName("font")
  if err != nil {
		log.Printf("Failed to get font static files with error[%s]", err.Error())
  }

  fontArr := make([]uriJSONObj, 0)
  for i := 0; i < len(fontFileNames); i++ {
    uriObj := &uriJSONObj{
      Uri: fmt.Sprintf(uriFormat, "font", fontFileNames[i]),
      Stamp: fmt.Sprintf("%d", i+1),
    }
    fontArr = append(fontArr, *uriObj);
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
