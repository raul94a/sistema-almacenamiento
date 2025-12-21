package power

import "net/http"

type Middleware func(http.Handler) http.Handler

func Handler(route string, handler http.HandlerFunc, middlewares ...Middleware){
    
	// Start with the final business logic handler
    var finalHandler http.Handler = handler

    // Wrap the handler in middlewares in reverse order
    // This ensures the first middleware in the slice is the first one executed
    for i := len(middlewares) - 1; i >= 0; i-- {
        finalHandler = middlewares[i](finalHandler)
    }

    // Register the route with the wrapped handler
    http.Handle(route, finalHandler)
}