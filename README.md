# MuxJWT: Go Module for JWT Auth
> MuxJWT: implement JWT authentication with gorilla/mux

## Installation
In a working directory with a ``go.mod`` file and [gorilla/mux](https://github.com/gorilla/mux) installed, type:
```
$ go get -u github.com/Tomer-Rubinstein/muxjwt
```

## Usage
Add the necessary imports:
```go
import (
  "github.com/Tomer-Rubinstein/muxjwt"
  "github.com/gorilla/mux"
  "net/http"
)
```
To initialize an authentication route:
```go
r := mux.NewRouter()
// Initialize a MuxJWT instace:
m := muxjwt.NewMuxJWT("my_secret", 60, "localhost")
// Create an "/auth" route that accepts POST requests:
m.InitAuthRoute(r, authFunc, "/auth", identifyFunc, "username", "password")
```
Where ``authFunc`` is the function responsible for authenticating given user credentials.

The "username" and "password" arguments are the POST request body data keys to pass their values to ``authFunc``.

And the ``identifyFunc`` returns the identifier of the request out of the given request body (in this case, the "username" field). It will be used as the subject claim in the payload of the JWT.

For simplicity sake, we'll just authenticate a single admin user, like so:
```go
func authFunc(body map[string]string) bool {
  username := body["username"]
  password := body["password"]
  return username == "admin" && password == "admin"
}
```
The ``identifyFunc`` as described:
```go
func identifyFunc(body map[string]string) string {
	return body["username"]
}
```
Now let's create 2 routes: one for the login page and the other be a secret page that only authenticated users can access.
```go
r.HandleFunc("/login", LoginHandler).Methods("GET")
m.ProtectedRoute(r, "/secret", SecretHandler).Methods("GET")
```
``ProtectedRoute`` creates a route that requires the request to satisy a cookie by the name ``"token_"+HOST`` with it's value being a JWT. If the token is valid, access to the set route("/secret") will be granted.

The handlers:
```go
func LoginHandler(w http.ResponseWriter, r *http.Request){
  http.ServeFile(w, r, "./static/LoginPage.html")
}

func SecretHandler(w http.ResponseWriter, r *http.Request){
  http.ServeFile(w, r, "./static/SecretPage.html")
}
```
``LoginPage.html`` contains the following form:
```HTML
<form action="/auth" method="post">
  <input type="text" id="username" name="username"><br>
  <input type="password" id="password" name="password"><br><br>
  <input type="submit" value="Submit">
</form>
```
The above issues (on submit) a POST request to "/auth" with the body data being ``username=<USER>&password=<PASSW>``.

``authFunc`` will then take these values and verify them.
If verified, MuxJWT will set a cookie named "token_localhost" with the value of the newly created JWT, on the client's browser.

**Complete sample code can be found at the "/example" directory.**