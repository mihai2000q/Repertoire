import { emptyPlaylist, reduxRender } from '../../test-utils.tsx'
import HomePlaylists from './HomePlaylists.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'
import { screen } from '@testing-library/react'

describe('Home Playlists', () => {
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

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<HomePlaylists />)

    expect(screen.getByText(/playlists/i)).toBeInTheDocument()
    for (const playlist of playlists) {
      expect(await screen.findByText(playlist.title)).toBeInTheDocument()
      if (playlist.imageUrl)
        expect(screen.getByRole('img', { name: playlist.title })).toHaveAttribute(
          'src',
          playlist.imageUrl
        )
      else
        expect(screen.getByLabelText(`default-icon-${playlist.title}`)).toBeInTheDocument()
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

    reduxRender(<HomePlaylists />)

    expect(await screen.findByText(/no playlists/i)).toBeInTheDocument()
  })
})
