package main

import (
	"fmt"
	"golang.org/x/net/netutil"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main(){
	HttpServer(Router())
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()


	mux.Handle("/img", http.HandlerFunc(Img))

	mux.Handle("/", http.HandlerFunc(Index))

	mux.Handle("/detail", http.HandlerFunc(Detail))

	mux.Handle("/friend", http.HandlerFunc(Friend))

	return mux
}

func Index(w http.ResponseWriter, r *http.Request) {
	str, _ := os.Getwd()

	if strings.HasPrefix(r.URL.Path,"/img"){
		log.Println("is img")
		file := str + "/img" +r.URL.Path[len("/str"):]
		log.Println(file)
		f,err := os.Open(file)
		defer f.Close()
		if err != nil && os.IsNotExist(err){
			file = str + "/default.jpg"
		}
		http.ServeFile(w,r,file)
		return
	}

	path := str + "/view/index.html"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_,_=fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	_,_=fmt.Fprintln(w, string(data))
	return
}

func Img(w http.ResponseWriter, r *http.Request) {
	str, _ := os.Getwd()

	if strings.HasPrefix(r.URL.Path,"/img"){
		file := str + r.URL.Path[len("/str"):]
		f,err := os.Open(file)
		defer f.Close()

		if err != nil && os.IsNotExist(err){
			file = str + "/default.jpg"
		}
		http.ServeFile(w,r,file)
		return
	}
}

func Detail(w http.ResponseWriter, r *http.Request) {
	str, _ := os.Getwd()
	path := str + "/view/detail.html"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_,_=fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	_,_=fmt.Fprintln(w, string(data))
	return
}

func Friend(w http.ResponseWriter, r *http.Request) {
	str, _ := os.Getwd()
	path := str + "/view/friend_list.html"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		w.WriteHeader(404)
		_,_=fmt.Fprintln(w, err)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	_,_=fmt.Fprintln(w, string(data))
	return
}

func HttpServer(router *http.ServeMux){
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := &http.Server{
		Addr:         ":28888",
		ReadTimeout:  4*time.Second,
		WriteTimeout: 4*time.Second,
		IdleTimeout:  4*time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:      router,
		// tls.Config ???????????? Certificates []Certificate
		// Certificate ???????????? Certificate PrivateKey ???????????? certFile keyFile ???????????????
	}

	// ????????????????????????????????????, ???????????????????????????????????????
	// ???????????? keep-alives ?????????????????? tcp ????????????????????? timewait ??????
	server.SetKeepAlivesEnabled(true)
	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic("Listen Err : "+ err.Error())
		return
	}
	defer l.Close()

	// ???????????????????????? ??????: linux/uinx???????????? win??????
	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println(err)
		return
	}


	log.Println("Starting http server port -> 28888")
	// ???????????????????????? ??????????????????????????? ????????????????????????
	// https://github.com/golang/net/blob/master/netutil/listen.go
	l = netutil.LimitListener(l, int(rLimit.Max)*10)
	err = server.Serve(l)
	if err != nil {
		panic("ListenAndServe err : "+ err.Error())
	}
}