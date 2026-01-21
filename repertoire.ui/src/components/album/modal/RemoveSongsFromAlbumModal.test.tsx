import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import RemoveSongsFromAlbumModal from './RemoveSongsFromAlbumModal.tsx'
import { RemoveSongsFromAlbumRequest } from '../../../types/requests/AlbumRequests.ts'
import { reduxRender, withToastify } from '../../../test-utils.tsx'

describe('Remove Songs From Album Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <RemoveSongsFromAlbumModal albumId={''} ids={['1', '2']} opened={true} onClose={vi.fn()} />
    )

    expect(screen.getByRole('dialog', { name: /remove songs from album/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /remove songs from album/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should remove songs from album and call onRemove', async () => {
    const user = userEvent.setup()

    let capturedRequest: RemoveSongsFromAlbumRequest
    server.use(
      http.put(`/albums/remove-songs`, async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromAlbumRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const albumId = 'album-id'
    const songIds = ['1', '2']
    const onClose = vi.fn()
    const onRemove = vi.fn()

    reduxRender(
      withToastify(
        <RemoveSongsFromAlbumModal
          albumId={albumId}
          ids={songIds}
          opened={true}
          onClose={onClose}
          onRemove={onRemove}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onRemove).toHaveBeenCalledOnce()
    expect(screen.getByText(`${songIds.length} songs removed from album!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: albumId, songIds: songIds })
  })
})
