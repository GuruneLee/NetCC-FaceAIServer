package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"errors"

	"github.com/Kagami/go-face"
	"github.com/gorilla/mux"
)

type Response struct {
	Feature face.Descriptor `json:"feature"`
	Message string          `json:"msg" `
	Err     string          `json:"error"`
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
			fmt.Errorf("getImg error, Error: ", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			var nilDes face.Descriptor
			failData := Response{nilDes, ErrMsg, err.Error()}
			failBytes, err := json.Marshal(failData)
			if err != nil {
				fmt.Errorf("marshaling error, Erro: ", err.Error())
				return
			}
			w.Write(failBytes)
			return
		}

		//get feature
		descriptor, err := getFeature(imgBin)
		if err != nil {
			fmt.Println("getFeature error, Erro: ", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			failData := Response{descriptor, ErrMsg, err.Error()}
			failBytes, err := json.Marshal(failData)
			if err != nil {
				fmt.Errorf("marshaling error, Erro: ", err.Error())
				return
			}
			w.Write(failBytes)
			return
		}
		featureData := Response{descriptor, SuccessMsg, ""}
		featureBytes, err := json.Marshal(featureData)
		if err != nil {
			fmt.Errorf("marshaling error, Erro: ", err.Error())
			return
		}
		//response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		w.Write(featureBytes)

	}).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getImg(r *http.Request) ([]byte, error) {
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("face-img")
	if err != nil {
		var emptyslice []byte
		return emptyslice, err
	}
	defer file.Close()

	imgByte, err := ioutil.ReadAll(file)
	if err != nil {
		var emptyslice []byte
		return emptyslice, err
	}
	return imgByte, err
}

func getFeature(bin []byte) (face.Descriptor, error) {
	modelDir := os.Getenv("MODEL_DIR")

	var nilDes face.Descriptor
	rec, err := face.NewRecognizer(modelDir)
	if err != nil {
		return nilDes, err
	}
	defer rec.Close()

	//this src can recognize only 'one-face'image
	f, err := rec.RecognizeSingle(bin)
	if f == nil {
		if err != nil {
			return nilDes, err
		} else {
			fmt.Errorf("There is no face on requested img\n")
			return nilDes, errors.New("There is no face on requested img")
		}
	}
	return f.Descriptor, err

}
