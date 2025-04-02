import { waitFor } from '@testing-library/react'

// Verify connect was called
import { Centrifuge } from 'centrifuge'
import { reduxRenderHook } from '../test-utils.tsx'
import { http, HttpResponse, ws } from 'msw'
import { setupServer } from 'msw/node'
import useCentrifuge from './useCentrifuge.ts'

describe('use Centrifuge', () => {
  const centrifugoUrl = 'wss://chat.example.com'

  const token = 'some-token'
  const handlers = [
    http.get('/centrifugo/token', async () => {
      return HttpResponse.json(token)
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => vi.stubEnv('VITE_CENTRIFUGO_URL', centrifugoUrl))

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => {
    vi.unstubAllEnvs()
    server.close()
  })

  it('should connect to Centrifugo and reuse existing connection', async () => {
    const connection = ws.link(centrifugoUrl)
    let numberOfConnections = 0
    server.use(
      connection.addEventListener('connection', () => {
        numberOfConnections++
        expect(numberOfConnections).toBe(1)
      })
    )

    const [{ result: firstResult }] = reduxRenderHook(() => useCentrifuge())

    expect(firstResult.current).toBeInstanceOf(Centrifuge)
    expect(firstResult.current.state).toBe(Centrifuge.State.Connecting)

    const [{ result: secondResult }] = reduxRenderHook(() => useCentrifuge())
    expect(secondResult.current).toBe(firstResult.current)
    await waitFor(() => expect(numberOfConnections).toBeGreaterThan(0))
  })

  it('should disconnect from Centrifugo when component unmounts', async () => {
    const connection = ws.link(centrifugoUrl)
    let connected = false
    server.use(
      connection.addEventListener('connection', () => {
        connected = true
      })
    )

    const [{ unmount, result }] = reduxRenderHook(() => useCentrifuge())

    await waitFor(() => expect(connected).toBeTruthy())

    unmount()

    await waitFor(() => {
      expect(result.current.state).toBe(Centrifuge.State.Disconnected)
    })
  })
})
