import { reduxRouterRender } from '../test-utils.tsx'
import RequireAuthentication from './RequireAuthentication.tsx'
import { Navigate, Route, Routes } from 'react-router-dom'
import { expect } from 'vitest'
import { screen } from '@testing-library/react'

describe('Require Authentication', () => {
  const render = (token: string | null) =>
    reduxRouterRender(
      <Routes>
        <Route element={<RequireAuthentication />}>
          <Route path={'/'} element={<div data-testid={'outlet'}>Outlet</div>} />
        </Route>

        <Route path={'sign-in'} element={<div data-testid={'sign-in'}>SignIn</div>} />
        <Route path={'*'} element={<Navigate to={'/'} replace />} />
      </Routes>,
      { auth: { token, historyOnSignIn: undefined } }
    )

  it('should remain on Outlet if user is authenticated', () => {
    render('my token')

    expect(window.location.pathname).toBe('/')
    expect(screen.getByTestId('outlet')).toBeInTheDocument()
  })

  it('should navigate to Sign In if user is not authenticated', () => {
    render(null)

    expect(window.location.pathname).toBe('/sign-in')
    expect(screen.getByTestId('sign-in')).toBeInTheDocument()
  })
})
