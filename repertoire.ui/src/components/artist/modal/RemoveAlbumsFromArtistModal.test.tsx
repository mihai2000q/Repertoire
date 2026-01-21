import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import RemoveAlbumsFromArtistModal from './RemoveAlbumsFromArtistModal.tsx'
import { RemoveAlbumsFromArtistRequest } from '../../../types/requests/ArtistRequests.ts'
import { reduxRender, withToastify } from '../../../test-utils.tsx'

describe('Remove Albums From Artist Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <RemoveAlbumsFromArtistModal artistId={''} ids={['1', '2']} opened={true} onClose={vi.fn()} />
    )

    expect(screen.getByRole('dialog', { name: /remove albums from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /remove albums from artist/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should remove albums from artist and call onRemove', async () => {
    const user = userEvent.setup()

    let capturedRequest: RemoveAlbumsFromArtistRequest
    server.use(
      http.put(`/artists/remove-albums`, async (req) => {
        capturedRequest = (await req.request.json()) as RemoveAlbumsFromArtistRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const artistId = 'artist-id'
    const albumIds = ['1', '2']
    const onClose = vi.fn()
    const onRemove = vi.fn()

    reduxRender(
      withToastify(
        <RemoveAlbumsFromArtistModal
          artistId={artistId}
          ids={albumIds}
          opened={true}
          onClose={onClose}
          onRemove={onRemove}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onRemove).toHaveBeenCalledOnce()
    expect(screen.getByText(`${albumIds.length} albums removed from artist!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: artistId, albumIds: albumIds })
  })
})
