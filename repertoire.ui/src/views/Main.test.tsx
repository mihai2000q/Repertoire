import { reduxRouterRender } from '../test-utils.tsx'
import Main from './Main.tsx'
import { screen } from '@testing-library/react'
import { Route, Routes } from 'react-router-dom'

describe('Main', () => {
  const render = (token?: string) =>
    reduxRouterRender(
      <Routes>
        <Route element={<Main />}>
          <Route path={'/'} element={<div data-testid={'outlet'}>Outlet</div>} />
        </Route>
      </Routes>,
      { auth: { token } }
    )

  it('should render and display just the outlet when user is not authenticated', () => {
    render()

    expect(screen.getByTestId('outlet')).toBeInTheDocument()
    expect(screen.queryByTestId('title-bar')).not.toBeInTheDocument() // env not desktop
  })

  it('should render and display just the outlet and title bar when user is not authenticated', () => {
    vi.stubEnv('VITE_PLATFORM', 'desktop')

    render()

    expect(screen.getByTestId('outlet')).toBeInTheDocument()
    expect(screen.getByTestId('title-bar')).toBeInTheDocument()
  })

  it('should render and display topbar, sidebar and outlet when user is authenticated', () => {
    render('some token')

    expect(screen.getByTestId('outlet')).toBeInTheDocument()
    expect(screen.getByRole('navigation')).toBeInTheDocument()
    expect(screen.getByRole('banner')).toBeInTheDocument()
  })

  it('should display title bar when platform is desktop', () => {
    vi.stubEnv('VITE_PLATFORM', 'desktop')

    render()

    expect(screen.getByTestId('title-bar')).toBeInTheDocument()
  })
})
