import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteSongsModal from './DeleteSongsModal.tsx'
import { BulkDeleteSongsRequest } from '../../../../types/requests/SongRequests.ts'

describe('Delete Songs Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<DeleteSongsModal ids={['1', '2']} opened={true} onClose={vi.fn()} />)

    expect(screen.getByRole('dialog', { name: /delete songs/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete songs/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete songs in bulk call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongsRequest
    server.use(
      http.put(`/songs/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteSongsModal ids={ids} opened={true} onClose={onClose} onDelete={onDelete} />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} songs deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ ids: ids })
  })
})
