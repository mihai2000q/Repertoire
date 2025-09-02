import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteArtistsModal from './DeleteArtistsModal.tsx'
import { BulkDeleteArtistsRequest } from '../../../../types/requests/ArtistRequests.ts'

describe('Delete Artists Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeleteArtistsModal ids={['1', '2']} opened={true} onClose={vi.fn()} />)

    expect(screen.getByRole('dialog', { name: /delete artists/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete artists/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /albums and songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete artists in bulk without associations and call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteArtistsRequest
    server.use(
      http.put(`/artists/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteArtistsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteArtistsModal ids={ids} opened={true} onClose={onClose} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} artists deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      ids: ids,
      withSongs: false,
      withAlbums: false
    })
  })

  it('should delete artist with associations', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteArtistsRequest
    server.use(
      http.put(`/artists/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteArtistsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteArtistsModal ids={ids} opened={true} onClose={onClose} />))

    await user.click(screen.getByRole('checkbox', { name: /albums and songs/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} artists deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      ids: ids,
      withSongs: true,
      withAlbums: true
    })
  })
})
