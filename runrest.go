package rest

/*
NEW WAY:
1) Let tbl := {table} from the request
2) Get tbl's TblDesc
3) Form the SQL stmt without a buffer ptr
4) It will use the TblDesc tools to return a buffer ptr as a RowModel (or "any")
5) Then can we print it as JSON, without ever really knowing its type ?

Perhaps each RowModel also needs to provide
its own JSON Decode & Encode functions.
Then everything occurs in text-only, without the code
ever knowing the specific types of the RowModels's ?
*/

import (
	L "github.com/fbaube/mlog"
	"net/http"
	"strconv"
//	S "strings"
//	"github.com/fbaube/m5db"
	"encoding/json"
	DRP "github.com/fbaube/datarepo"
//	DRM "github.com/fbaube/datarepo/rowmodels"
//	DRS "github.com/fbaube/datarepo/sqlite"
)

var sRestPortNr string

var the_m5db DRP.SimpleRepo

// THIS IS JUST A DATASTORE !
/*
type taskServer struct {
	store *taskstore.TaskStore
}
func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}
*/

func RunRest(portNr int, aDB DRP.SimpleRepo) error {
	if portNr == 0 { // env.RestPort
		return nil
	}
	if aDB == nil { return nil }
	println("runrest.go: RUN REST")
	sRestPortNr = strconv.Itoa(portNr)
	the_m5db = aDB 
	
// 2023.10 https://eli.thegreenplace.net/2023/better-http-server-routing-in-go-122
// 2021.01 https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library

	mux := http.NewServeMux()
/*
	mux.HandleFunc(  "POST /{table}/",      hCreate)
	mux.HandleFunc("DELETE /{table}/",      hDeleteAll)
	mux.HandleFunc("DELETE /{table}/{id}/", hDeleteByID)
	mux.HandleFunc(   "GET /tag/{tag}/",    hMultiGetByTag)
*/
	mux.HandleFunc(   "GET /{table}/",      hGetAll)
	mux.HandleFunc(   "GET /{table}/{id}/", hGetByID)
/*
	huma.Get(api, "/about",
		pRsp.ContentType = "text/html"
		sTop, eTop := GetContents("top.hf")
		sMid, eMid := GetContents("mid.hf")
		sAbt, eAbt := GetContents("about.hb")
		sBtm, eBtm := GetContents("btm.hf")
		pRsp.Body = []byte(
			// "<!DOCTYPE html>\n<html>\n<body>\nABOUT!\n" +
			// "</body></html>")
			sTop + sMid + sAbt + sBtm)
		return pRsp, cmp.Or(eTop, eMid, eAbt, eBtm, nil)
*/
	// fmt.Printf("API: %+v \n", api)
	println("==> Running REST server on port:", 8000) // sRestPortNr)
	// Start the server!
	http.ListenAndServe("127.0.0.1:8000", mux)
	// http.ListenAndServe("localhost:8888", mux)
	return nil
}

// TIME ABOUT HEALTH CONTACT RSS FUNC CFG ENV ADMIN STC DB

/*
	// ADMIN
	r.HandleFunc("/stc", hdlStcRoot)
	// TOPICS, MAPS, DATABASE, STATIC CONTENT
	rtrTpc := r.PathPrefix("/tpc").Subrouter()
	rtrMap := r.PathPrefix("/map").Subrouter()
	rtrDb := r.PathPrefix("/db").Subrouter()
	// HOME (incl. "About", etc.)
	r.HandleFunc("/", HomeHandler)

	// TOPICS
	rtrTpc.HandleFunc("/{id}/meta", hdlTopicMeta)
	rtrTpc.HandleFunc("/{id}/links", hdlTopicLinks)
	rtrTpc.HandleFunc("/{id}", hdlTopic)
	rtrTpc.HandleFunc("/", hdlTopicRoot)

	// MAPS
	rtrMap.HandleFunc("/{id}/meta", hdlMapMeta)
	rtrMap.HandleFunc("/{id}/links", hdlMapLinks)
	rtrMap.HandleFunc("/{id}", hdlMap)
	rtrMap.HandleFunc("/", hdlMapRoot)

	// DB (schemas? stats?)
	rtrDb.HandleFunc("/fld/{name}", hdlDbField)
	rtrDb.HandleFunc("/{name}", hdlDbTable)
	rtrDb.HandleFunc("/", hdlDbRoot)

	// go func() {
	if err := http.ListenAndServe(":"+sRestPortNr, r); err != nil {
		L.L.Error("REST server failed: " + err.Error())
	}
	return nil
}

* /

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("home vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	* /
	L.L.Info(s)
}

func TopicRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("topic vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)'

		fmt.Fprintf(w, s)
	* /
	L.L.Info(s)
}
func MapRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("map vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	* /
	L.L.Info(s)
}

func StcRootHdlr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var s string
	s = r.RequestURI + ": " + fmt.Sprintf("static vars: %+v\n", vars)
	/*
		println(s)
		ssnLog.Println(s)
		fmt.Fprintf(w, s)
	* /
	L.L.Info(s)

	// This will serve files under http://localhost:8000/static/<filename>
	// r.PathPrefix("/s/").Handler(http.StripPrefix("/s/", http.FileServer(http.Dir(dir))))
}

*/

// func DRS.DoSelectByIdGeneric[T DRM.RowModel](pSR *SqliteRepo, anID int, pDest T) (bool, error) 

func hGetByID(w http.ResponseWriter, req *http.Request) {
  L.L.Info("Handling get-by-id at %s", req.URL.Path)

  tbl := req.PathValue("table")
  if tbl == "" {
    http.Error(w, "invalid table", http.StatusBadRequest)
    return
  }
  id, err := strconv.Atoi(req.PathValue("id"))
  if err != nil {
    http.Error(w, "invalid id", http.StatusBadRequest)
    return
  }
  // ok, e := DRS.DoSelectByIdGeneric(the_m5db, id, typedBuf)
  rq := new(DRP.RestQuery)
  rq.DBOp = http.MethodGet // "get"
  rq.Table = tbl
  rq.Id1 = id
  // RunQuery1(*DRU.QuerySpec) (any, error)
  // Any, e := RunQuery1(rq)
  
  if err != nil {
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }

  // renderJSON(w, task)
}

func hGetAll(w http.ResponseWriter, req *http.Request) {
  L.L.Info("Handling get-all at %s", req.URL.Path)

  tbl := req.PathValue("table")
  if tbl == "" {
    http.Error(w, "invalid table", http.StatusBadRequest)
    return
  }
  // ok, e := DRS.DoSelectByIdGeneric(the_m5db, id, typedBuf)
  rq := new(DRP.RestQuery)
  rq.DBOp = "get"
  rq.Table = tbl
  rq.Id1 = -1
/*
  if err != nil {
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
*/
  // renderJSON(w, task)
}

// func DRS.DoInsertGeneric[T DRM.RowModel](pSR *SqliteRepo, pRM T) (int, error) 

/*
The only handler that's a bit special is hCreate, since it
has to parse JSON data sent by the client in the request body.
There are some nuances to JSON parsing in requests that I did
not cover - check out this post for a more thorough approach.
https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

*/

// Hyper-robust: How to Parse a JSON Request Body in Go
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
  js, err := json.Marshal(v)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

