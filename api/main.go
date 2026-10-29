package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{"The Go Programming Language", "Donovan & Kernighan"},
	{"Designing Data-Intensive Applications", "Martin Kleppmann"},
	{"The Pragmatic Programmer", "Thomas & Hunt"},
}

type Claims struct {
	Username    string `json:"preferred_username"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func main() {
	keys := oidc.NewRemoteKeySet(context.Background(), os.Getenv("JWKS_URL"))
	verifier := oidc.NewVerifier(os.Getenv("ISSUER"), keys, &oidc.Config{ClientID: os.Getenv("AUDIENCE")})

	mux := http.NewServeMux()
	mux.Handle("GET /api/books", requireRole(verifier, "reader", listBooks))
	log.Println("api listening on :8181")
	log.Fatal(http.ListenAndServe(":8181", cors(os.Getenv("WEB_ORIGIN"), mux)))
}

func requireRole(v *oidc.IDTokenVerifier, role string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		// checks signature, issuer, audience and expiry
		token, err := v.Verify(r.Context(), raw)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		var c Claims
		if err := token.Claims(&c); err != nil {
			http.Error(w, "bad claims", http.StatusUnauthorized)
			return
		}
		if !slices.Contains(c.RealmAccess.Roles, role) {
			http.Error(w, c.Username+" does not have the "+role+" role", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func cors(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
