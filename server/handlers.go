package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

type merchantPutRequest struct {
	Link string `json:"link" binding:"required"`
}

type errorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID, err := strconv.Atoi(vars["id"])
	if err != nil {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Couldn't convert merchant ID to int"}); err != nil {
			panic(err)
		}
		return
	}
	body := &merchantPutRequest{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Code: http.StatusUnprocessableEntity,
			Msg:  "Couldn't decode body to struct"}); err != nil {
			panic(err)
		}
		return
	}
	url, err := url.ParseRequestURI(body.Link)
	if err != nil {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Code: http.StatusUnprocessableEntity,
			Msg:  "Couldn't parse link"}); err != nil {
			panic(err)
		}
		return
	}
	taskID := s.Data.Cache.CreateTask()

	w.Header().Set("Content-Location", fmt.Sprintf("../tasks/%d", taskID))
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(errorResponse{
		Code: http.StatusAccepted,
		Msg:  fmt.Sprintf("Task with ID=%d created at Content-Location", taskID)}); err != nil {
		panic(err)
	}
	go s.download(*url, taskID, merchantID)
}

func (s *Server) getTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Couldn't convert task ID to int"}); err != nil {
			panic(err)
		}
		return
	}
	stats, err := s.Data.Cache.GetTaskByID(id)
	if err != nil {
		if err := json.NewEncoder(w).Encode(errorResponse{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("Couldn't find task with given ID=%d", id)}); err != nil {
			panic(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) getOffers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var oID, mID int
	var sub string
	var err error
	if _, ok := vars["offerID"]; ok {
		oID, err = strconv.Atoi(vars["offerID"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errorResponse{
				Code: http.StatusBadRequest,
				Msg:  "Couldn't convert offerID to int"}); err != nil {
				panic(err)
			}
			return
		}
	}
	if _, ok := vars["merchantID"]; ok {
		mID, err = strconv.Atoi(vars["merchantID"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errorResponse{
				Code: http.StatusBadRequest,
				Msg:  "Couldn't convert merchantID to int"}); err != nil {
				panic(err)
			}
			return
		}
	}
	if _, ok := vars["sub"]; ok {
		sub = vars["sub"]
	}
	merchants := s.Data.GetMerchants(oID, mID, sub)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(merchants)
}
