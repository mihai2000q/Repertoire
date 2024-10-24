import { reduxRouterRender } from '../test-utils.tsx'
import Main from './Main.tsx'
import { screen } from '@testing-library/react'
import { Route, Routes } from 'react-router-dom'

describe('Main', () => {
  const render = () =>
    reduxRouterRender(
      <Routes>
        <Route element={<Main />}>
          <Route path={'/'} element={<div>Outlet</div>} />
        </Route>
      </Routes>
    )

  it('should render and display outlet, sidebar and topbar', () => {
    render()

    expect(screen.getByText('Outlet')).toBeInTheDocument()
    expect(screen.getByRole('navigation')).toBeInTheDocument()
    expect(screen.getByRole('banner')).toBeInTheDocument()
  })
})
