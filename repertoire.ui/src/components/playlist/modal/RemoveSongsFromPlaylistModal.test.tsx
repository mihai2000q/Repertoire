import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import RemoveSongsFromPlaylistModal from './RemoveSongsFromPlaylistModal.tsx'
import { RemoveSongsFromPlaylistRequest } from '../../../types/requests/PlaylistRequests.ts'
import { reduxRender, withToastify } from '../../../test-utils.tsx'

describe('Remove Songs From Playlist Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <RemoveSongsFromPlaylistModal
        playlistId={''}
        ids={['1', '2']}
        opened={true}
        onClose={vi.fn()}
      />
    )

    expect(screen.getByRole('dialog', { name: /remove songs from playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /remove songs from playlist/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should remove songs from playlist and call onRemove', async () => {
    const user = userEvent.setup()

    let capturedRequest: RemoveSongsFromPlaylistRequest
    server.use(
      http.put(`/playlists/songs/remove`, async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromPlaylistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const playlistId = 'playlist-id'
    const playlistSongIds = ['1', '2']
    const onClose = vi.fn()
    const onRemove = vi.fn()

    reduxRender(
      withToastify(
        <RemoveSongsFromPlaylistModal
          playlistId={playlistId}
          ids={playlistSongIds}
          opened={true}
          onClose={onClose}
          onRemove={onRemove}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onRemove).toHaveBeenCalledOnce()
    expect(
      screen.getByText(`${playlistSongIds.length} songs removed from playlist!`)
    ).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: playlistId, playlistSongIds: playlistSongIds })
  })
})
