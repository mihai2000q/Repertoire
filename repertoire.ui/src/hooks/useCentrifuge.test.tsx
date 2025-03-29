import { waitFor } from '@testing-library/react'
import useCentrifuge from './useCentrifuge.ts'

// Verify connect was called
import { Centrifuge } from 'centrifuge'
import { reduxRenderHook } from '../test-utils.tsx'
import { http, HttpResponse, ws } from 'msw'
import { setupServer } from 'msw/node'

describe('use Centrifuge', () => {
  const centrifugoUrl = 'wss://chat.example.com'

  const token = 'some-token'
  const handlers = [
    http.get('/auth/centrifugo', async () => {
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
    let numberOfConnections = 0
    const connection = ws.link(centrifugoUrl)

    const [{ result: firstResult }] = reduxRenderHook(() => useCentrifuge())

    expect(firstResult.current).toBeInstanceOf(Centrifuge)
    expect(firstResult.current.state).toBe(Centrifuge.State.Connecting)
    connection.addEventListener('connection', () => {
      expect(firstResult.current.state).toBe(Centrifuge.State.Connected)
      numberOfConnections++
      expect(numberOfConnections).toBe(1)
    })

    const [{ result: secondResult }] = reduxRenderHook(() => useCentrifuge())
    expect(secondResult.current).toBe(firstResult.current)
  })

  it('should disconnect from Centrifugo when component unmounts', async () => {
    const connection = ws.link(centrifugoUrl)

    const [{ unmount, result }] = reduxRenderHook(() => useCentrifuge())

    connection.addEventListener('connection', () => {
      expect(result.current.state).toBe(Centrifuge.State.Connected)
    })

    unmount()

    await waitFor(() => expect(result.current.state).toBe(Centrifuge.State.Disconnected))
  })
})
