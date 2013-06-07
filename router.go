package nutty

import (
  "strings"
  "net/http"
)

type ControllerWithIndex interface {
  Index(*App, http.ResponseWriter, *http.Request)
}

type ControllerWithCreate interface {
  Create(*App, http.ResponseWriter, *http.Request)
}

type ControllerWithDestroy interface {
  Destroy(*App, http.ResponseWriter, *http.Request)
}

type ControllerWithUpdate interface {
  Update(*App, http.ResponseWriter, *http.Request)
}

type Router struct {
  handlers          map[string](map[string]interface{})
  initializations   map[string]bool
}


func (routes *Router) Map(uri string, controller interface{}, httpMethods []string, nuttyApp *App) {
  if !routes.initializations[uri] {
    routes.initializations[uri] = true
    routes.handlers[uri] = make(map[string]interface{})
    http.HandleFunc(uri, func(resp http.ResponseWriter, req *http.Request) {
      if routes.handlers[uri][req.Method] == nil {
        http.NotFound(resp, req)
      } else {
        if req.Method == "POST" {
          if req.FormValue("_method") == "delete" || req.FormValue("_method") == "DELETE" {
            (routes.handlers[uri][req.Method]).(ControllerWithDestroy).Destroy(nuttyApp, resp, req)
          } else if req.FormValue("_method") == "put" || req.FormValue("_method") == "PUT" {
            (routes.handlers[uri][req.Method]).(ControllerWithUpdate).Update(nuttyApp, resp, req)
          } else {
            (routes.handlers[uri][req.Method]).(ControllerWithCreate).Create(nuttyApp, resp, req)
          }
        } else {
          (routes.handlers[uri][req.Method]).(ControllerWithIndex).Index(nuttyApp, resp, req)
        }
      }
    })
  }
  
  for _, method := range httpMethods {
    routes.handlers[uri][strings.ToUpper(method)] = controller
  }
}

// Defaults to GET if no http methods sent
func (routes *Router) Root(ctrl ControllerWithIndex, nuttyApp *App) {
  if !routes.initializations["/"] {
    routes.initializations["/"] = true
    routes.handlers["/"] = make(map[string]interface{})
  }
  
  handler := func(w http.ResponseWriter, r *http.Request) { ctrl.Index(nuttyApp,w,r) }
  routes.handlers["/"]["GET"] = handler
  http.HandleFunc("/", handler)
  http.HandleFunc("/index", handler)
  http.HandleFunc("/index.html", handler)
}
