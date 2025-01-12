import { reduxRouterRender } from '../test-utils.tsx'
import Main from './Main.tsx'
import { screen } from '@testing-library/react'
import { Route, Routes } from 'react-router-dom'

describe('Main', () => {
  const render = (token?: string) =>
    reduxRouterRender(
      <Routes>
        <Route element={<Main />}>
          <Route path={'/'} element={<div>Outlet</div>} />
        </Route>
      </Routes>,
      { auth: { token } }
    )

  it('should render and display just the outlet when user is not authenticated', () => {
    render()

    expect(screen.getByText('Outlet')).toBeInTheDocument()
    expect(screen.queryByTestId('title-bar')).not.toBeInTheDocument() // env not desktop
  })

  it('should render and display just the outlet and title bar when user is not authenticated', () => {
    import.meta.env.VITE_PLATFORM = 'desktop'

    render()

    expect(screen.getByText('Outlet')).toBeInTheDocument()
    expect(screen.getByTestId('title-bar')).toBeInTheDocument()

    // restore
    import.meta.env.VITE_PLATFORM = ''
  })

  it('should render and display topbar, sidebar and outlet when user is authenticated', () => {
    render('some token')

    expect(screen.getByText('Outlet')).toBeInTheDocument()
    expect(screen.getByRole('navigation')).toBeVisible()
    expect(screen.getByRole('banner')).toBeVisible()
  })

  it('should display title bar when platform is desktop', () => {
    import.meta.env.VITE_PLATFORM = 'desktop'

    render()

    expect(screen.getByTestId('title-bar')).toBeVisible()

    // restore
    import.meta.env.VITE_PLATFORM = ''
  })
})
