package main

import(
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseForm(); err != nil{
		fmt.Fprintf(w,"ParseForm() err : %v", err)
	}
	fmt.Fprintf(w,"Post request successfull")
	first_name := r.FormValue("fname")
	last_name := r.FormValue("lname")

	fmt.Fprintf(w,"\nHi %v %v",first_name,last_name)

}

func helloHandler(w http.ResponseWriter, r *http.Request)  {
	if r.URL.Path != "/hello"{
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	
	if r.Method != "GET"{
		http.Error(w, "Method not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w,"Hey This Is RK")
}

func main(){
	fileserver := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileserver)
	http.HandleFunc("/form",formHandler)
	http.HandleFunc("/hello",helloHandler)

	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080",nil); err != nil{
		log.Fatal(err)
	}
}