import { AuthProvider } from './context/AuthContext'

export default function App() {
  return (
    <AuthProvider>
      <div>Red Packet App</div>
    </AuthProvider>
  )
}
