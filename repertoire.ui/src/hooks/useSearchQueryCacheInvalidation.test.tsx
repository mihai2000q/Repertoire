import { http, HttpResponse, ws } from 'msw'
import { setupServer } from 'msw/node'
import { reduxRenderHook } from '../test-utils.tsx'
import useSearchQueryCacheInvalidation from './useSearchQueryCacheInvalidation.ts'
import {expect} from "vitest";
import useCentrifuge from "./useCentrifuge.ts";
import {waitFor} from "@testing-library/react";

describe('use Search Query Cache Invalidation', () => {
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

  it('should subscribe to search channel', async () => {
    const search = ws.link(centrifugoUrl)
    const userId = 'some-user-id'

    let connected = false
    server.use(
      search.addEventListener('connection', () => {
        expect(result.current.getSubscription(`search:${userId}`)).toBeDefined()
        connected = true
      })
    )

    const [{ result }] = reduxRenderHook(() => useCentrifuge())

    reduxRenderHook(() => useSearchQueryCacheInvalidation(), {
      global: {
        userId: userId,
        artistDrawer: undefined,
        albumDrawer: undefined,
        songDrawer: undefined
      }
    })
    await waitFor(() => expect(connected).toBeTruthy())
  })
})
