// - When your proxy receives an HTTP request for an object from a browser,
// it generates a new HTTP request for the same object and sends it to the origin server.
//
// - When the proxy receives the corresponding HTTP response with the object from the origin server,
// it creates a new HTTP response, including the object, and sends it to the client.
//
// - Handle simultaneous connections by using multiple threads

package main
