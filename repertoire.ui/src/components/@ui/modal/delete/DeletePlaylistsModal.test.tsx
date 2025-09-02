import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeletePlaylistsModal from './DeletePlaylistsModal.tsx'
import { BulkDeletePlaylistsRequest } from '../../../../types/requests/PlaylistRequests.ts'

describe('Delete Playlists Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeletePlaylistsModal ids={['1', '2']} opened={true} onClose={vi.fn()} />)

    expect(screen.getByRole('dialog', { name: /delete playlists/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete playlists/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete playlists in bulk call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeletePlaylistsRequest
    server.use(
      http.put(`/playlists/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeletePlaylistsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeletePlaylistsModal ids={ids} opened={true} onClose={onClose} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} playlists deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ ids: ids })
  })
})
