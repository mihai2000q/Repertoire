import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteAlbumsModal from './DeleteAlbumsModal.tsx'
import { BulkDeleteAlbumsRequest } from '../../../../types/requests/AlbumRequests.ts'

describe('Delete Albums Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeleteAlbumsModal ids={['1', '2']} opened={true} onClose={vi.fn()} />)

    expect(screen.getByRole('dialog', { name: /delete albums/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete albums/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /songs/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete albums in bulk without songs and call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteAlbumsRequest
    server.use(
      http.put(`/albums/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteAlbumsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteAlbumsModal ids={ids} opened={true} onClose={onClose} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} albums deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      ids: ids,
      withSongs: false
    })
  })

  it('should delete album with songs', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteAlbumsRequest
    server.use(
      http.put(`/albums/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteAlbumsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()

    reduxRender(withToastify(<DeleteAlbumsModal ids={ids} opened={true} onClose={onClose} />))

    await user.click(screen.getByRole('checkbox', { name: /songs/i }))
    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} albums deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      ids: ids,
      withSongs: true
    })
  })
})
