import { Navigate, Outlet } from 'react-router-dom'
import useAuth from '../hooks/useAuth'
import { ReactElement } from 'react'

function IsAlreadyAuthenticated(): ReactElement {
  const isAuthenticated = useAuth()
  return isAuthenticated ? <Navigate to={'home'} replace /> : <Outlet />
}

export default IsAlreadyAuthenticated
