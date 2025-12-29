Custom Handlers

In the previous exercise, we used the http.FileServer function, which simply returns a built-in http.Handler.

An http.Handler is just an interface:

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

Any type with a ServeHTTP method that matches the http.HandlerFunc signature above is an http.Handler. Take a second to think about it: it makes a lot of sense! To handle an incoming HTTP request, all a function needs is a way to write a response and the request itself.
Assignment

Let's add a readiness endpoint to the Chirpy server! Readiness endpoints are commonly used by external systems to check if our server is ready to receive traffic.

The endpoint should be accessible at the /healthz path using any HTTP method.

The endpoint should simply return a 200 OK status code indicating that it has started up successfully and is listening for traffic. The endpoint should return a Content-Type: text/plain; charset=utf-8 header, and the body will contain a message that simply says "OK" (the text associated with the 200 status code).

Later this endpoint can be enhanced to return a 503 Service Unavailable status code if the server is not ready.
1. Add the Readiness Endpoint

I recommend using the mux.HandleFunc to register your handler. Your handler can just be a function that matches the signature of http.HandlerFunc:

func(http.ResponseWriter, *http.Request)

Your handler should do the following:

    Write the Content-Type header
    Write the status code using w.WriteHeader
    Write the body text using w.Write

2. Update the Fileserver Path

Now that we've added a new handler, we don't want potential conflicts with the fileserver handler. Update the fileserver to use the /app/ path instead of /.

Not only will you need to mux.Handle the /app/ path, you'll also need to strip the /app prefix from the request path before passing it to the fileserver handler. You can do this using the http.StripPrefix function.

Run and submit the CLI tests.
Assignment

Let's add a readiness endpoint to the Chirpy server! Readiness endpoints are commonly used by external systems to check if our server is ready to receive traffic.

The endpoint should be accessible at the /healthz path using any HTTP method.

The endpoint should simply return a 200 OK status code indicating that it has started up successfully and is listening for traffic. The endpoint should return a Content-Type: text/plain; charset=utf-8 header, and the body will contain a message that simply says "OK" (the text associated with the 200 status code).

Later this endpoint can be enhanced to return a 503 Service Unavailable status code if the server is not ready.
1. Add the Readiness Endpoint

I recommend using the mux.HandleFunc to register your handler. Your handler can just be a function that matches the signature of http.HandlerFunc:

func(http.ResponseWriter, *http.Request)

Your handler should do the following:

    Write the Content-Type header
    Write the status code using w.WriteHeader
    Write the body text using w.Write

2. Update the Fileserver Path

Now that we've added a new handler, we don't want potential conflicts with the fileserver handler. Update the fileserver to use the /app/ path instead of /.

Not only will you need to mux.Handle the /app/ path, you'll also need to strip the /app prefix from the request path before passing it to the fileserver handler. You can do this using the http.StripPrefix function.

Run and submit the CLI tests.
Assignment

Let's add a readiness endpoint to the Chirpy server! Readiness endpoints are commonly used by external systems to check if our server is ready to receive traffic.

The endpoint should be accessible at the /healthz path using any HTTP method.

The endpoint should simply return a 200 OK status code indicating that it has started up successfully and is listening for traffic. The endpoint should return a Content-Type: text/plain; charset=utf-8 header, and the body will contain a message that simply says "OK" (the text associated with the 200 status code).

Later this endpoint can be enhanced to return a 503 Service Unavailable status code if the server is not ready.
1. Add the Readiness Endpoint

I recommend using the mux.HandleFunc to register your handler. Your handler can just be a function that matches the signature of http.HandlerFunc:

func(http.ResponseWriter, *http.Request)

Your handler should do the following:

    Write the Content-Type header
    Write the status code using w.WriteHeader
    Write the body text using w.Write

2. Update the Fileserver Path

Now that we've added a new handler, we don't want potential conflicts with the fileserver handler. Update the fileserver to use the /app/ path instead of /.

Not only will you need to mux.Handle the /app/ path, you'll also need to strip the /app prefix from the request path before passing it to the fileserver handler. You can do this using the http.StripPrefix function.

Run and submit the CLI tests.
