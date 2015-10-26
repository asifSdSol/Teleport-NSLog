package main

import (
    "net/http"
    "net/http/httputil"
    "io/ioutil"
    "compress/gzip"
    "io"
    "os" 
    "github.com/kardianos/osext"
    "path"
    "bytes"
    "log"
)

func Handler(resp http.ResponseWriter, req *http.Request) {
    query := req.URL.Query()
    deviceid := query.Get("devid")
    log.Println("Device id: ", deviceid)
    if deviceid == "" {

        d, err := httputil.DumpRequest(req, true)
        if err != nil {
            log.Println(err.Error())
        } else {
            log.Printf("ERROR: Missing devid. \nRequest:\n%s", d)
        }
        http.Error(resp, "devid is required parameter", http.StatusBadRequest)
        return

    }

    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        http.Error(resp, err.Error(), http.StatusInternalServerError)
        return
    }   
    
    r, err := gzip.NewReader(bytes.NewReader(body))
    if err != nil {
        http.Error(resp, err.Error(), http.StatusInternalServerError)
        return
    }   
    execpath,execpatherr := osext.Executable()
    if execpatherr != nil {
	http.Error(resp, execpatherr.Error(), http.StatusInternalServerError)
        log.Printf("Error when locating executable path: %s", execpatherr.Error())
        return
    }
    execdirpath := path.Dir(execpath)
    logfilepath := path.Join(execdirpath,"logs",deviceid + ".log")
    dst, err := os.OpenFile(logfilepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
    log.Printf("Appending to %s", logfilepath)
    if err != nil {
        http.Error(resp, err.Error(), http.StatusInternalServerError)
	log.Printf("Error when opening file: %s", err.Error())
        return
    }
    defer dst.Close()
    defer r.Close()
    io.Copy(dst, r)
}

func main() {
    http.HandleFunc("/", Handler)
    http.ListenAndServe(":8081", nil)
}
