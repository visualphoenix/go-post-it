package main

import (
	"crypto/md5"
	"fmt"
	"io"
        "mime/multipart"  
	"net/http"
	"path/filepath" 
	"os"
	"strconv"
	"text/template"
	"time"
)

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
        var (  
             status int  
             err  error  
        )  
        defer func() {  
             if nil != err {  
                  http.Error(w, err.Error(), status)  
             }  
        }() 
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else if r.Method == "PUT" {
		err = os.MkdirAll("./test/"+filepath.Dir(r.URL.Path[1:]), 0755)
		if err != nil {
			fmt.Fprint(w, "Fail", err)
			fmt.Println("[Notice]", "Put Fail", r.URL.Path, r.RemoteAddr, err)
		}
		f, err := os.OpenFile("./test/"+r.URL.Path[1:], os.O_WRONLY|os.O_CREATE, 0666)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                defer f.Close()

		_, err = io.Copy(f, r.Body)
		if err != nil {
			fmt.Fprint(w, "Fail", err)
			fmt.Println("[Notice]", "Put Fail", r.URL.Path, r.RemoteAddr, err)
			return
		}
		fmt.Println("[Info]", "Put", r.URL.Path)
		fmt.Fprint(w, "Success")
	} else {
          // parse request  
          const _24K = (1 << 20) * 24  
          err = r.ParseMultipartForm(_24K);
          if err != nil {  
               status = http.StatusInternalServerError  
               return  
          }  
          for _, fheaders := range r.MultipartForm.File {  
               for _, hdr := range fheaders {  
                    // open uploaded  
                    var infile multipart.File  
                    if infile, err = hdr.Open(); nil != err {  
                         status = http.StatusInternalServerError  
                         return  
                    }  
                    // open destination  
                    var outfile *os.File  
                    if outfile, err = os.OpenFile("./test/" + hdr.Filename, os.O_WRONLY|os.O_CREATE, 0666); nil != err {  
                         status = http.StatusInternalServerError  
                         return  
                    }  
                    // 32K buffer copy  
                    var written int64  
                    if written, err = io.Copy(outfile, infile); nil != err {  
                         status = http.StatusInternalServerError  
                         return  
                    }  
                    w.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written)) + "\n"))  
               }  
          }
	}
}

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(upload))
}

