Goroutines in Servers

In Go, goroutines are used to serve many requests at the same time, but not all servers are quite so performant.

Go was built by Google, and one of the purposes of its creation was to power Google's massive web infrastructure. Go's goroutines are a great fit for web servers because they're lighter weight than operating system threads, but still take advantage of multiple cores. Let's compare a Go web server's concurrency model to other popular languages and frameworks.
Node.js / Express.js

In JavaScript land, servers are typically single-threaded. A Node.js server (often using the Express framework) only uses one CPU core at a time. It can still handle many requests at once by using an async event loop. That just means whenever a request has to wait on I/O (like to a database), the server puts it on pause and does something else for a bit.

This might sound horribly inefficient, but it's not too bad. Node servers do just fine with the I/O workloads associated with most CRUD apps (Where processing is offloaded to the Database). You only start to run into trouble with this model when you need your server to do CPU-intensive work.
Node.js / Express.js

In JavaScript land, servers are typically single-threaded. A Node.js server (often using the Express framework) only uses one CPU core at a time. It can still handle many requests at once by using an async event loop. That just means whenever a request has to wait on I/O (like to a database), the server puts it on pause and does something else for a bit.

This might sound horribly inefficient, but it's not too bad. Node servers do just fine with the I/O workloads associated with most CRUD apps (Where processing is offloaded to the Database). You only start to run into trouble with this model when you need your server to do CPU-intensive work.
Takeaways

    Go servers are great for performance whether the workload is I/O or CPU-bound
    Node.js and Express work well for I/O-bound tasks, but struggle with CPU-bound tasks

I'm not saying Go is always "better" than JavaScript when it comes to back-end development, but it generally outperforms when it comes to computational speed.

Click to hide video
Your browser does not support playing HTML5 video. You can instead. Here is a description of the content: go vs js servers

Takeaways

    Go servers are great for performance whether the workload is I/O or CPU-bound
    Node.js and Express work well for I/O-bound tasks, but struggle with CPU-bound tasks

I'm not saying Go is always "better" than JavaScript when it comes to back-end development, but it generally outperforms when it comes to computational speed.

Click to hide video
Your browser does not support playing HTML5 video. You can instead. Here is a description of the content: go vs js servers

