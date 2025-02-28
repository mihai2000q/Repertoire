import { emptyUser, reduxRouterRender } from '../../test-utils.tsx'
import Topbar from './Topbar.tsx'
import { AppShell } from '@mantine/core'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'

describe('Topbar', () => {
  const render = (token: string | null = 'some token', toggleSidebar: () => void = () => {}) =>
    reduxRouterRender(
      <AppShell>
        <Topbar toggleSidebar={toggleSidebar} />
      </AppShell>,
      { auth: { token } }
    )

  const handlers = [
    http.get('/users/current', async () => {
      return HttpResponse.json(emptyUser)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should display just the search bar, when token is not available', async () => {
    server.use(
      http.get('/users/current', async () => {
        return HttpResponse.json(undefined, { status: 404 })
      })
    )

    render(undefined)

    expect(screen.getByRole('searchbox', { name: 'search' })).toBeInTheDocument()
  })

  it('should display search bar and user, when token is available', async () => {
    render('some token')

    expect(screen.getByRole('searchbox', { name: 'search' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'user' })).toBeInTheDocument()
  })

  it('should display navigation buttons when on desktop, when token is available', () => {
    vi.stubEnv('VITE_PLATFORM', 'desktop')

    render()

    expect(screen.getByRole('button', { name: 'back' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'forward' })).toBeInTheDocument()
  })

  it.skip('should display sidebar menu button, when screen is small', async () => {
    const userEventDispatcher = userEvent.setup()

    const originalInnerWidth = window.innerWidth;
    vi.spyOn(window, 'innerWidth', 'get').mockImplementation(() => 500);

    const toggleSidebar = vitest.fn()

    render('some token', toggleSidebar)

    const button = screen.getByRole('button', { name: 'toggle-sidebar' })
    expect(button).toBeInTheDocument()
    await userEventDispatcher.click(button)
    expect(toggleSidebar).toHaveBeenCalledOnce()

    // restore
    vi.spyOn(window, 'innerWidth', 'get').mockImplementation(() => originalInnerWidth);
  })
})
