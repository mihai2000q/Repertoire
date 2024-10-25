import { reduxRouterRender } from '../test-utils.tsx'
import IsAlreadyAuthenticated from './IsAlreadyAuthenticated.tsx'
import {Navigate, Route, Routes} from 'react-router-dom'
import { expect } from 'vitest'
import { screen } from '@testing-library/react'

describe('Is Already Authenticated', () => {
  const render = (token: string | null) =>
    reduxRouterRender(
      <Routes>
        <Route element={<IsAlreadyAuthenticated />}>
          <Route path={'/'} element={<div>Outlet</div>} />
        </Route>

        <Route path={'/home'} element={<div>Home</div>} />
        <Route path={'*'} element={<Navigate to={'/'} replace />} />
      </Routes>,
      { auth: { token } }
    )

  it('should remain on Outlet if user is not authenticated', () => {
    render(null)

    expect(window.location.pathname).toBe('/')
    expect(screen.getByText('Outlet')).toBeInTheDocument()
  })

  it('should navigate to Home if user is already authenticated', () => {
    render('my token')

    expect(window.location.pathname).toBe('/home')
    expect(screen.getByText('Home')).toBeInTheDocument()
  })
})