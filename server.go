package main

import (
  "github.com/gorilla/pat"
  "net/http"
  "log"
  "encoding/json"
  "strconv"
  "io/ioutil"
  "time"
)

type HttpResponse struct {
  url string
  body string
  err error
}

func WriteJSONResponse(responses map[string]string, code int, res http.ResponseWriter) {
  js, err := json.MarshalIndent(responses, "", "  ")
  if err != nil {
    log.Printf("JSON serialization error: %v", err)
    return
  }
  res.Header().Set("Content-Length", strconv.Itoa(len(js)+1))
  res.WriteHeader(code)
  res.Write(js)
  res.Write([]byte("\n"))
  return
}

func MultigetHandler(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-Type", "text/javascript")
  urls := req.URL.Query()["url"]
  ch := make(chan *HttpResponse)
  for _, url := range urls {
    go func(url string) {
      resp, err := http.Get(url)
      if err != nil {
        log.Printf("Encountered an error with %s: %v", url, err)
        ch <- &HttpResponse{url, "", nil}
      } else {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        ch <- &HttpResponse{url, string(body), err}
      }
    }(url)
  }
  responses := make(map[string]string)
  for {
    select {
      case r := <-ch:
        log.Printf("Got %s", r.url)
        responses[r.url] = string(r.body)
        if len(responses) == len(urls) {
          WriteJSONResponse(responses, 200, res)
          return
        }
      case <-time.After(5 * time.Second):
        if len(responses) > 0 {
          WriteJSONResponse(responses, 203, res)
        } else {
          WriteJSONResponse(responses, 204, res)
        }
        return
    }
  }
}

func main() {
  r := pat.New()
  r.Get("/multiget", MultigetHandler)
  http.Handle("/", r)
  err := http.ListenAndServe(":12345", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
