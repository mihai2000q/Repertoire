import { emptyUser, reduxRouterRender } from '../test-utils.tsx'
import Main from './Main.tsx'
import { screen } from '@testing-library/react'
import { Route, Routes } from 'react-router-dom'
import { setupServer } from 'msw/node'
import { http, HttpResponse, ws } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { SearchBase } from '../types/models/Search.ts'
import { beforeEach } from 'vitest'

describe('Main', () => {
  const handlers = [
    http.get('/centrifugo/token', async () => {
      return HttpResponse.json('')
    }),
    http.get('/users/current', async () => {
      return HttpResponse.json(emptyUser)
    }),
    http.get('/search', () => {
      const response: WithTotalCountResponse<SearchBase> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => {
    server.listen()
    vi.mock('../hooks/useMainScroll.ts', () => ({
      default: vi.fn(() => ({
        ref: { current: document.createElement('div') }
      }))
    }))
  })

  beforeEach(() => {
    const centrifugoUrl = 'wss://chat.example.com'
    vi.stubEnv('VITE_CENTRIFUGO_URL', centrifugoUrl)
    const search = ws.link(centrifugoUrl)
    server.use(search.addEventListener('connection', () => {}))
  })

  afterEach(() => server.resetHandlers())

  afterAll(() => {
    vi.unstubAllEnvs()
    vi.clearAllMocks()
    server.close()
  })

  const render = (token?: string) =>
    reduxRouterRender(
      <Routes>
        <Route element={<Main />}>
          <Route path={'/'} element={<div data-testid={'outlet'}>Outlet</div>} />
        </Route>
      </Routes>,
      { auth: { token, historyOnSignIn: { index: 0, justSignedIn: false } } }
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
