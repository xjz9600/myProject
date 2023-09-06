package web

type Middleware func(HandleFunc) HandleFunc
