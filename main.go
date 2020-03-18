package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	AddResource(router, new(HelloResource))

	log.Fatal(http.ListenAndServe(":8120", router))
}

type HelloResource struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (HelloResource) Uri() string {
	return "/hello"
}

func (HelloResource) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{200, "message", map[string]interface{}{
		"key1": "value1",
	}}
}

type ContentsWrite struct {
	GetNotSupported
	PutNotSupported
	DeleteNotSupported
}

func (ContentsWrite) Uri() string {
	return "/ContentsWrite"
}

type Document struct {
	DocumentName string `json:"documentName"`
	DocumentId   string `json:"documentId"`
	Writer       string `json:"writer"`
}

var Documents []Document

func (ContentsWrite) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{200, "", Documents}
}

func (ContentsWrite) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	Documents := r.FormValue("data")
	if len(Documents) > 0 {
		// 로직 구현
		return Response{200, "success", nil}
	} else {
		return Response{400, "fail", nil}
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Resource interface {
	Uri() string
	Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
	Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response
}

type (
	GetNotSupported    struct{}
	PostNotSupported   struct{}
	PutNotSupported    struct{}
	DeleteNotSupported struct{}
)

func (GetNotSupported) Get(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PostNotSupported) Post(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (PutNotSupported) Put(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func (DeleteNotSupported) Delete(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) Response {
	return Response{405, "", nil}
}

func abort(rw http.ResponseWriter, statusCode int) {
	rw.WriteHeader(statusCode)
}

func HttpResponse(rw http.ResponseWriter, req *http.Request, res Response) {
	content, err := json.Marshal(res)

	if err != nil {
		abort(rw, 500)
	}

	rw.WriteHeader(res.Code)
	rw.Write(content)
}

func AddResource(router *httprouter.Router, resource Resource) {
	fmt.Println("\"" + resource.Uri() + "\" api 등록")

	router.GET(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Get(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.POST(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Post(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.PUT(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Put(rw, r, ps)
		HttpResponse(rw, r, res)
	})
	router.DELETE(resource.Uri(), func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		res := resource.Delete(rw, r, ps)
		HttpResponse(rw, r, res)
	})
}
