package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"serialdemo/service"
)

var publicKey, _ = service.LoadPublicKey("public.key")

func scale(w http.ResponseWriter, r *http.Request) {
	var info service.WeightInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if info.Error != "" {
		log.Printf("received error: %+v", info.Error)
		//TODO need to retry
		return
	} else if !verifySignature(info) {
		log.Println("invalid signature!")
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	// Do something with the WeightInfo struct...
	log.Printf("WeightInfo: %+v", info.WeightInfoToSign)
}

func verifySignature(info service.WeightInfo) bool {
	jsonBytes, _ := json.Marshal(info.WeightInfoToSign)
	r, s := new(big.Int).SetBytes(info.R), new(big.Int).SetBytes(info.S)
	return ecdsa.Verify(publicKey, service.Hash(jsonBytes), r, s)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/scale", scale)
	url := "0.0.0.0:8080"
	log.Println("app backend listening:", url)
	if err := http.ListenAndServe(url, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
