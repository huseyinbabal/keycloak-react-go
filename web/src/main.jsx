import { createRoot } from 'react-dom/client'
import { keycloak } from './keycloak'
import App from './App'
import './style.css'

await keycloak.init({ onLoad: 'check-sso', pkceMethod: 'S256' })

createRoot(document.getElementById('root')).render(<App />)
