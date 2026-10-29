import { useState } from 'react'
import { keycloak } from './keycloak'

export default function App() {
  const [books, setBooks] = useState([])
  const [error, setError] = useState('')

  async function loadBooks() {
    setError('')
    // refresh the access token if it expires in the next 30 seconds
    await keycloak.updateToken(30)
    const res = await fetch('http://localhost:8181/api/books', {
      headers: { Authorization: `Bearer ${keycloak.token}` },
    })
    if (!res.ok) {
      setBooks([])
      setError(`${res.status} ${await res.text()}`)
      return
    }
    setBooks(await res.json())
  }

  if (!keycloak.authenticated) {
    return (
      <main>
        <h1>Bookstore</h1>
        <button onClick={() => keycloak.login()}>Sign in</button>
      </main>
    )
  }

  const { iss, aud, realm_access, exp } = keycloak.tokenParsed

  return (
    <main>
      <h1>Bookstore</h1>
      <p>Signed in as <b>{keycloak.tokenParsed.preferred_username}</b></p>
      <button onClick={loadBooks}>Load books</button>
      <button onClick={() => keycloak.logout()}>Sign out</button>
      <div className="cols">
        <section>
          {error && <p className="error">{error}</p>}
          <ul>
            {books.map((b) => (
              <li key={b.title}>{b.title} <span>{b.author}</span></li>
            ))}
          </ul>
        </section>
        <pre>{JSON.stringify({ iss, aud, roles: realm_access.roles, exp: new Date(exp * 1000).toLocaleTimeString() }, null, 2)}</pre>
      </div>
    </main>
  )
}
