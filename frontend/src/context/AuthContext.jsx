import { createContext, useReducer } from 'react'

export const AuthContext = createContext(null)

function authReducer(state, action) {
  switch (action.type) {
    case 'LOGIN':
      return { token: action.token, isAuthenticated: true }
    case 'LOGOUT':
      return { token: null, isAuthenticated: false }
    default:
      return state
  }
}

function getInitialState() {
  const token = localStorage.getItem('token')
  return {
    token,
    isAuthenticated: !!token,
  }
}

export function AuthProvider({ children }) {
  const [state, dispatch] = useReducer(authReducer, null, getInitialState)

  function login(token) {
    localStorage.setItem('token', token)
    dispatch({ type: 'LOGIN', token })
  }

  function logout() {
    localStorage.removeItem('token')
    dispatch({ type: 'LOGOUT' })
  }

  return (
    <AuthContext.Provider value={{ ...state, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
