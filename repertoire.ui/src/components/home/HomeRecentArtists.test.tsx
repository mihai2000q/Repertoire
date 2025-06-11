import { emptyArtist, reduxRouterRender } from '../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Artist from '../../types/models/Artist.ts'
import { screen } from '@testing-library/react'
import HomeRecentArtists from './HomeRecentArtists.tsx'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../state/store.ts'
import { expect } from 'vitest'

describe('Home Recent Artists', () => {
  const artists: Artist[] = [
    {
      ...emptyArtist,
      id: '1',
      name: 'Artist 1'
    },
    {
      ...emptyArtist,
      id: '2',
      name: 'Artist 2',
      imageUrl: 'something.png'
    }
  ]

  const handlers = [
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: artists,
        totalCount: artists.length
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
    reduxRouterRender(<HomeRecentArtists />)

    expect(screen.getByText(/recent artists/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'forward' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'back' })).toBeInTheDocument()
    expect(screen.getByTestId('artists-loader')).toBeInTheDocument()
    for (const artist of artists) {
      expect(await screen.findByText(artist.name)).toBeInTheDocument()
      if (artist.imageUrl)
        expect(screen.getByRole('img', { name: artist.name })).toHaveAttribute(
          'src',
          artist.imageUrl
        )
      else expect(screen.getByLabelText(`default-icon-${artist.name}`)).toBeInTheDocument()
    }
  })

  it('should display empty message when there are no artists', async () => {
    server.use(
      http.get('/artists', async () => {
        const response: WithTotalCountResponse<Artist> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRouterRender(<HomeRecentArtists />)

    expect(await screen.findByText(/no artists/i)).toBeInTheDocument()
  })

  it('should open artist drawer on click', async () => {
    const user = userEvent.setup()

    const artist = artists[1]

    const [_, store] = reduxRouterRender(<HomeRecentArtists />)

    await user.click(await screen.findByRole('img', { name: artist.name }))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(artist.id)
  })

  it('should display menu on artist right click', async () => {
    const user = userEvent.setup()

    const artist = artists[0]

    reduxRouterRender(<HomeRecentArtists />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByLabelText(`default-icon-${artist.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to artist when clicking on view details', async () => {
      const user = userEvent.setup()

      const artist = artists[0]

      reduxRouterRender(<HomeRecentArtists />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByLabelText(`default-icon-${artist.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/artist/${artist.id}`)
    })
  })
})
