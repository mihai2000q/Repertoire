import Artist from 'src/types/models/Artist.ts'
import {
  defaultSongFiltersMetadata,
  emptyArtist,
  reduxRouterRender,
  withToastify
} from '../../test-utils.tsx'
import ArtistCard from './ArtistCard.tsx'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { RootState } from '../../state/store.ts'

describe('Artist Card', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const handlers = [
    http.get('/playlists', async () => {
      return HttpResponse.json(defaultSongFiltersMetadata)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render with minimal info', () => {
    reduxRouterRender(<ArtistCard artist={artist} />)

    expect(screen.getByLabelText(`default-icon-${artist.name}`)).toBeInTheDocument()
    expect(screen.getByText(artist.name)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localArtist: Artist = {
      ...artist,
      imageUrl: 'something.png'
    }

    reduxRouterRender(<ArtistCard artist={localArtist} />)

    expect(screen.getByRole('img', { name: localArtist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localArtist.name })).toHaveAttribute(
      'src',
      localArtist.imageUrl
    )
    expect(screen.getByText(localArtist.name)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<ArtistCard artist={artist} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${artist.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /open drawer/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open artist drawer when clicking on open drawer', async () => {
      const user = userEvent.setup()

      const [_, store] = reduxRouterRender(withToastify(<ArtistCard artist={artist} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${artist.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /open drawer/i }))

      expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
      expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(artist.id)
    })

    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<ArtistCard artist={artist} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${artist.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<ArtistCard artist={artist} />)

    await user.click(screen.getByLabelText(`default-icon-${artist.name}`))
    expect(window.location.pathname).toBe(`/artist/${artist.id}`)
  })
})
