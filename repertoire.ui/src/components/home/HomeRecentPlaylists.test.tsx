import { emptyPlaylist, reduxRouterRender } from '../../test-utils.tsx'
import HomeRecentPlaylists from './HomeRecentPlaylists.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import { RootState } from '../../state/store.ts'
import { userEvent } from '@testing-library/user-event'

describe('Home Recent Playlists', () => {
  const playlists: Playlist[] = [
    {
      ...emptyPlaylist,
      id: '1',
      title: 'Playlist 1'
    },
    {
      ...emptyPlaylist,
      id: '2',
      title: 'Playlist 2',
      imageUrl: 'something.png'
    }
  ]

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = {
        models: playlists,
        totalCount: playlists.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRouterRender(<HomeRecentPlaylists />)

    expect(screen.getByText(/recent playlists/i)).toBeInTheDocument()
    for (const playlist of playlists) {
      expect(await screen.findByText(playlist.title)).toBeInTheDocument()
      if (playlist.imageUrl)
        expect(screen.getByRole('img', { name: playlist.title })).toHaveAttribute(
          'src',
          playlist.imageUrl
        )
      else expect(screen.getByLabelText(`default-icon-${playlist.title}`)).toBeInTheDocument()
    }
  })

  it('should display empty message when there are no playlists', async () => {
    server.use(
      http.get('/playlists', async () => {
        const response: WithTotalCountResponse<Playlist> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeRecentPlaylists />)

    expect(await screen.findByText(/no playlists/i)).toBeInTheDocument()
  })

  it('should open drawer on playlist click', async () => {
    const user = userEvent.setup()

    const playlist = playlists[1]

    const [, store] = reduxRouterRender(<HomeRecentPlaylists />)

    await user.click(await screen.findByRole('img', { name: playlist.title }))

    expect((store.getState() as RootState).global.playlistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.playlistDrawer.playlistId).toBe(playlist.id)
  })

  it('should display menu on playlist right click', async () => {
    const user = userEvent.setup()

    const playlist = playlists[0]

    reduxRouterRender(<HomeRecentPlaylists />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByLabelText(`default-icon-${playlist.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to playlist when clicking on view details', async () => {
      const user = userEvent.setup()

      const playlist = playlists[0]

      reduxRouterRender(<HomeRecentPlaylists />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByLabelText(`default-icon-${playlist.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/playlist/${playlist.id}`)
    })
  })
})
