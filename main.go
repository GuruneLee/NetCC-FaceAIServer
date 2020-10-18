package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	goFace "github.com/Kagami/go-face"
	"github.com/gorilla/mux"
)

type Response struct {
	Feature goFace.Descriptor `json:"feature"`
	Message string            `json:"msg" `
	Err     error             `json:"error"`
}

const (
	SuccessMsg string = "feature responsed"
	ErrMsg     string = "you got some errors"
)

func main() {
	fmt.Println("server start")
	router := mux.NewRouter()

	router.HandleFunc("/get/feature", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("client entering")

		//parsing
		imgBin, err := getImg(r)
		if err != nil {
			fmt.Println("getImg error, Erro: ", fmt.Sprint(err))
			return
		}

		//get feature
		descriptor, err := getFeature(imgBin)
		if err != nil {
			fmt.Println("getFeature error, Erro: ", fmt.Sprint(err))
			return
		}
		featureData := Response{descriptor, SuccessMsg, nil}
		featureBytes, _ := json.Marshal(featureData)

		//response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(featureBytes)

	}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getImg(r *http.Request) ([]byte, error) {
	file, _, err := r.FormFile("face-img")
	if err != nil {
		var emptyslice []byte
		return emptyslice, err
	}
	defer file.Close()
	imgByte, err := ioutil.ReadAll(file)
	return imgByte, nil
}

func getFeature(bin []byte) (goFace.Descriptor, error) {
	modleDir := os.Getenv("MODEL_DIR")

	var nilDes goFace.Descriptor
	rec, err := goFace.NewRecognizer(modleDir)
	if err != nil {
		return nilDes, err
	}
	face, err := rec.RecognizeSingle(bin)
	if err != nil {
		return nilDes, err
	}

	return face.Descriptor, nil

}
