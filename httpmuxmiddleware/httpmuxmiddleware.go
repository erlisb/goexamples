package main

import (
	"flag"
	"github.com/NYTimes/gziphandler"
	"github.com/evolutiontechnologies/gorillaerp/module/ioxer"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type (
	MyHandlerV3 func(http.ResponseWriter, *http.Request, *ioxer.IOXer)
)

func homeHTML() string {
	return `
		<html>
			<head>
				<script src="https://code.jquery.com/jquery-2.x-git.min.js"></script>
				<script src="static/debug.js"></script>
			</head>
			<body>
				<select id="User">
					<option value="user1">User1</option>
					<option value="user2">User2</option>
					<option value="user3">User3</option>
				</select>
				<button id="SalOrderBtn">SalOrder</button>
				<button id="NewsBtn">News</button>
			</body>
		</html>
`
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Host: %v", r.Host)
	// r.URL.Scheme will be empty if you're accessing the HTTP server not from an HTTP proxy,
	// a browser can issue a relative HTTP request instead of a absolute URL.
	// Additionally, you could check in the server/handler whether you get a
	// relative or absolute URL in the request by calling the IsAbs() method.
	// Reference: http://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server
	log.Printf("Scheme: %v", r.URL.Scheme)
	log.Printf("IsAbs: %v", r.URL.IsAbs())

	if _, err := w.Write([]byte(homeHTML())); err != nil {
		log.Printf("w.Write: %v", err.Error())
	}
}

func srvNotFound(w http.ResponseWriter, r *http.Request, o *ioxer.IOXer) {
	o.AddError("404 page not found")
}

func srvUserAuthentication(w http.ResponseWriter, r *http.Request, o *ioxer.IOXer) {
	Authorization := r.Header.Get("Authorization")

	log.Printf("Authorization: %v", Authorization)

	switch Authorization {
	case "user2":
		o.AddError("You do not have the permission")
	case "user3":
		o.AddError("You do not have the permission")
	}
}

func srvSalOrder(w http.ResponseWriter, r *http.Request, o *ioxer.IOXer) {
	so := &struct {
		DealerName string
		IDShipAddr int64
		Price      float64
	}{}

	i := ioxer.NewIOXer()
	i.PutObj(so)

	if err := i.Decode(r.Body); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("SQLDealerName: %s", i.GetData("SQLDealerName"))
	log.Printf("SQLIDShipAddr: %d", int(i.GetData("SQLIDShipAddr").(float64)))
	log.Printf("SQLPrice: %f", i.GetData("SQLPrice"))
	log.Printf("JSON: %s", i.EncodePretty())

	vars := mux.Vars(r)
	IDSalOrder := vars["IDSalOrder"]

	o.PutData("IDSalOrder", IDSalOrder)
}

func srvNews(w http.ResponseWriter, r *http.Request, o *ioxer.IOXer) {
	news := struct {
		Subject string
		Author  string
		Body    string
	}{
		Subject: "Hello World",
		Author:  "Jun",
		Body:    "This is a Hello World message",
	}

	o.PutObj(news)

}

func srvNoChecking(mh ...MyHandlerV3) http.HandlerFunc {
	return handlerLoopV3(append([]MyHandlerV3{}, mh...))
}

func srvRegularChecking(mh ...MyHandlerV3) http.HandlerFunc {
	return handlerLoopV3(append([]MyHandlerV3{srvUserAuthentication}, mh...))
}

func handlerLoopV3(myhandlers []MyHandlerV3) http.HandlerFunc {
	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := ioxer.NewIOXer()

		for _, myhandler := range myhandlers {
			if myhandler(w, r, o); o.ErrCount > 0 {
				o.Echo(w)
				return
			}
		}

		o.Echo(w)
	})).(http.HandlerFunc)
}

func main() {
	//
	var dir string

	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the static dir under the current dir")
	flag.Parse()

	//
	r := mux.NewRouter()

	// 404 not found handler
	//r.NotFoundHandler = http.HandlerFunc(srvNotFound)
	r.NotFoundHandler = srvNoChecking(srvNotFound)

	//
	r.HandleFunc("/", srvHome)

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	// Only matches if domain is "www.example.com".
	s := r.Host("erp.local").Subrouter()

	s.HandleFunc("/SalOrder/{IDSalOrder}", srvRegularChecking(srvSalOrder))

	//
	s2 := r.Host("erp.local").Subrouter()
	s2.HandleFunc("/News", srvNoChecking(srvNews))

	//
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
	}
}
