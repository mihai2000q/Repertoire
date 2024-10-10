import { Navigate, Outlet, useLocation } from 'react-router-dom'
import useAuth from '../hooks/useAuth'
import { ReactElement } from 'react'

function RequireAuthentication(): ReactElement {
  const location = useLocation()
  const isAuthenticated = useAuth()

  return isAuthenticated ? (
    <Outlet />
  ) : (
    <Navigate to={'sign-in'} state={{ from: location }} replace />
  )
}

export default RequireAuthentication
