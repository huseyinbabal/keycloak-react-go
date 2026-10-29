# Keycloak with React and Go

A bookstore with three parts, all started by one `docker compose up`:

- **Keycloak 26.8** (+ Postgres 18) holds the users and signs the tokens
- **web**: a React app that signs in with keycloak-js (authorization code flow + PKCE S256)
- **api**: a Go API that checks the access token with go-oidc and needs the `reader` realm role

```bash
docker compose up -d
```

| Service | URL |
|---|---|
| Keycloak admin console | http://localhost:8180 (admin / admin) |
| React app | http://localhost:5173 |
| Go API | http://localhost:8181/api/books |

The `bookstore` realm is imported on start from `realm/bookstore-realm.json` (client `web`, audience mapper `api`, role `reader`, two users). Keycloak generates new signing keys on the first start.

| User | Password | Role | `/api/books` |
|---|---|---|---|
| alice | alice | reader | 200, the book list |
| bob | bob | none | 403 `bob does not have the reader role` |

```bash
curl -i localhost:8181/api/books                                              # 401 missing token
curl -i -H "Authorization: Bearer made.up.token" localhost:8181/api/books     # 401 invalid token
```

All passwords here are for a local machine only. Keycloak runs with `start-dev` (HTTP, relaxed hostname checks); production needs `start` with a hostname and TLS.

Write-up: https://huseyinbabal.dev/posts/keycloak-react-go-pkce-jwt-roles/
Video: https://youtu.be/4hc0AmtjoOI
