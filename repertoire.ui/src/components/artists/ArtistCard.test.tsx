import Artist from 'src/types/models/Artist.ts'
import { reduxRouterRender, withToastify } from '../../test-utils.tsx'
import ArtistCard from './ArtistCard.tsx'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

describe('Artist Card', () => {
  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRouterRender(<ArtistCard artist={artist} />)

    expect(screen.getByRole('img', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByText(artist.name)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<ArtistCard artist={artist} />)

    await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: artist.name }) })
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      reduxRouterRender(withToastify(<ArtistCard artist={artist} />))

      // Assert
      await user.pointer({ keys: '[MouseRight>]', target: screen.getByRole('img', { name: artist.name }) })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<ArtistCard artist={artist} />)

    await user.click(screen.getByRole('img', { name: artist.name }))
    expect(window.location.pathname).toBe(`/artist/${artist.id}`)
  })
})
