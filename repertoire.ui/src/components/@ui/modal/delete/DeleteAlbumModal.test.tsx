import { http, HttpResponse } from 'msw'
import { emptyAlbum, reduxRender, withToastify } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import Album from '../../../../types/models/Album.ts'
import { setupServer } from 'msw/node'
import DeleteAlbumModal from './DeleteAlbumModal.tsx'
import { userEvent } from '@testing-library/user-event'

describe('Delete Album Modal', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1'
  }

  const handlers = [
    http.delete(`/albums/${album.id}`, () => {
      return HttpResponse.json({ message: 'it worked' })
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeleteAlbumModal opened={true} onClose={vi.fn()} album={album} />)

    expect(screen.getByRole('dialog', { name: /delete album/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete album and call onDelete', async () => {
    const user = userEvent.setup()

    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteAlbumModal opened={true} onClose={onClose} album={album} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
  })

  it('should delete album without associations', async () => {
    const user = userEvent.setup()

    let withSongs: string | null
    server.use(
      http.delete(`/albums/${album.id}`, (req) => {
        const params = new URL(req.request.url).searchParams
        withSongs = params.get('withSongs')
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteAlbumModal opened={true} onClose={onClose} album={album} />))

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(withSongs).toBe('false')
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
  })

  it('should delete album with associations', async () => {
    const user = userEvent.setup()

    let withSongs: string | null
    server.use(
      http.delete(`/albums/${album.id}`, (req) => {
        const params = new URL(req.request.url).searchParams
        withSongs = params.get('withSongs')
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteAlbumModal opened={true} onClose={onClose} album={album} />))

    await user.click(screen.getByRole('checkbox', { name: /songs/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(withSongs).toBe('true')
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
  })
})
