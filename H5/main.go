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
		// tls.Config 有个属性 Certificates []Certificate
		// Certificate 里有属性 Certificate PrivateKey 分别保存 certFile keyFile 证书的内容
	}

	// 如果在高频高并发的场景下, 有很多请求是可以复用的时候
	// 最好开启 keep-alives 减少三次握手 tcp 销毁连接时有个 timewait 时间
	server.SetKeepAlivesEnabled(true)
	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic("Listen Err : "+ err.Error())
		return
	}
	defer l.Close()

	// 开启最高连接数， 注意: linux/uinx有效果， win无效
	var rLimit syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println(err)
		return
	}


	log.Println("Starting http server port -> 28888")
	// 对连接数的保护， 设置为最高连接数是 本机的最高连接数
	// https://github.com/golang/net/blob/master/netutil/listen.go
	l = netutil.LimitListener(l, int(rLimit.Max)*10)
	err = server.Serve(l)
	if err != nil {
		panic("ListenAndServe err : "+ err.Error())
	}
}