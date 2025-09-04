import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import RemoveSongsFromArtistModal from './RemoveSongsFromArtistModal.tsx'
import { RemoveSongsFromArtistRequest } from '../../../types/requests/ArtistRequests.ts'
import { reduxRender, withToastify } from '../../../test-utils.tsx'

describe('Remove Songs From Artist Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <RemoveSongsFromArtistModal artistId={''} ids={['1', '2']} opened={true} onClose={vi.fn()} />
    )

    expect(screen.getByRole('dialog', { name: /remove songs from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /remove songs from artist/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should remove songs from artist and call onRemove', async () => {
    const user = userEvent.setup()

    let capturedRequest: RemoveSongsFromArtistRequest
    server.use(
      http.put(`/artists/remove-songs`, async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const artistId = 'artist-id'
    const songIds = ['1', '2']
    const onClose = vi.fn()
    const onRemove = vi.fn()

    reduxRender(
      withToastify(
        <RemoveSongsFromArtistModal
          artistId={artistId}
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
    expect(screen.getByText(`${songIds.length} songs removed from artist!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: artistId, songIds: songIds })
  })
})
