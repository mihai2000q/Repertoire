import { http, HttpResponse } from 'msw'
import { emptyArtist, reduxRender, withToastify } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import Artist from '../../../types/models/Artist.ts'
import { setupServer } from 'msw/node'
import DeleteArtistModal from './DeleteArtistModal.tsx'
import { userEvent } from '@testing-library/user-event'

describe('Delete Artist Modal', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const handlers = [
    http.delete(`/artists/${artist.id}`, () => {
      return HttpResponse.json({ message: 'it worked' })
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeleteArtistModal opened={true} onClose={vi.fn()} artist={artist} />)

    expect(screen.getByRole('dialog', { name: /delete artist/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete artist/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /albums and songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete artist and call onDelete', async () => {
    const user = userEvent.setup()

    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteArtistModal opened={true} onClose={onClose} artist={artist} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
  })

  it('should delete artist without associations', async () => {
    const user = userEvent.setup()

    let withAlbums: string | null
    let withSongs: string | null
    server.use(
      http.delete(`/artists/${artist.id}`, (req) => {
        const params = new URL(req.request.url).searchParams
        withAlbums = params.get('withAlbums')
        withSongs = params.get('withSongs')
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteArtistModal opened={true} onClose={onClose} artist={artist} />))

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(withAlbums).toBe('false')
    expect(withSongs).toBe('false')
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
  })

  it('should delete artist with associations', async () => {
    const user = userEvent.setup()

    let withAlbums: string | null
    let withSongs: string | null
    server.use(
      http.delete(`/artists/${artist.id}`, (req) => {
        const params = new URL(req.request.url).searchParams
        withAlbums = params.get('withAlbums')
        withSongs = params.get('withSongs')
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteArtistModal opened={true} onClose={onClose} artist={artist} />))

    await user.click(screen.getByRole('checkbox', { name: /albums and songs/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(withAlbums).toBe('true')
    expect(withSongs).toBe('true')
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
  })
})
