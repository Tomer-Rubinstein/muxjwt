# MuxJWT: Go Module for JWT Auth
> MuxJWT: implement JWT authentication with gorilla/mux

## Installation
In a working directory with a ``go.mod``  file and [gorilla/mux](https://github.com/gorilla/mux) installed, type:
```go
$ go get -u github.com/Tomer-Rubinstein/muxjwt
```

## Usage
Add the necessary imports:
```go
import (
	"github.com/Tomer-Rubinstein/muxjwt"
	"github.com/gorilla/mux"
)
```
To initialize an authentication route:
```go
r := mux.NewRouter()
// The following creates an "/auth" route that accepts
// only POST requests.
muxjwt.InitAuthRoute(r, authFunc, "/auth", "username", "pasword")
```
Where ``authFunc`` is the function responsible for authenticating given user credentials. I.e. lookup user in a database and check for a matching password.
The "username" and "password" arguments are the POST request body data keys to pass their values to ``authFunc``.

For simplicity sake, we'll just authenticate a single admin user, like so:
```go
func authFunc(body map[string]string) bool {
	username := body["username"]
	password := body["password"]
	return username == "admin" && password == "admin"
}
```
Now let's create 2 routes: one for the login page and the other be a secret page that only authenticated users can access.
```go
r.HandleFunc("/login", LoginHandler).Methods("GET")
muxjwt.ProtectedRoute(r, "/secret", SecretHandler).Methods("GET")
```
Using the ``ProtectedRoute`` function, MuxJWT will require the request to satisfy the ``Authorization`` header with the value of: ``Bearer <token>``.

The handlers:
```go
func LoginHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/LoginPage.html")
}

func SecretHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./static/SecretPage.html")
}
```
LoginPage.html contains the following form:
```HTML
<form action="/auth" method="post">
	<input type="text" id="username" name="username"><br>
	<input type="password" id="password" name="password"><br><br>
	<input type="submit" value="Submit">
</form>
```
Say we submitted ``username=admin&password=admin`` as our POST body data to ``/auth``, MuxJWT will then verify the data and return a JWT accordingly.
**This JWT should be stored in client's localStorage**

Now, to access the ``/secret`` page, we shall pass ``Authorization: Bearer <jwt>`` as a header to our GET request and if the token is valid and still not expired, we will get the contents of the secret page!
