import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteSongPartsModal from './DeleteSongPartsModal.tsx'
import { BulkDeleteSongPartsRequest } from '../../types/requests/SongPartRequests.ts'

describe('Delete Song Parts Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <DeleteSongPartsModal ids={['1', '2']} songId={'1'} opened={true} onClose={vi.fn()} />
    )

    expect(screen.getByRole('dialog', { name: /delete parts/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete parts/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should close when cancelling', async () => {
    const user = userEvent.setup()
    const onClose = vi.fn()

    reduxRender(
      <DeleteSongPartsModal ids={['1', '2']} songId={'song-1'} opened={true} onClose={onClose} />
    )

    await user.click(screen.getByRole('button', { name: /cancel/i }))

    expect(onClose).toHaveBeenCalledOnce()
  })

  it('should explain which duplicate parts will be ignored', () => {
    reduxRender(
      <DeleteSongPartsModal
        ids={['1', '1', '2', '2', '2']}
        songId={'1'}
        opened={true}
        onClose={vi.fn()}
      />
    )

    expect(
      screen.getByText(
        'Are you sure you want to delete 2 parts? (3 duplicates will be ignored)'
      )
    ).toBeInTheDocument()
  })

  it('should delete song parts in bulk call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongPartsRequest
    server.use(
      http.put(`/songs/parts/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongPartsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const songId = '1'
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteSongPartsModal
          ids={ids}
          songId={songId}
          opened={true}
          onClose={onClose}
          onDelete={onDelete}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
    expect(screen.getByText(`${ids.length} parts deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ ids: ids, songId: songId })
  })

  it('should delete each unique part only once', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongPartsRequest
    server.use(
      http.put(`/songs/parts/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongPartsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '1', '2', '2', '2']
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteSongPartsModal
          ids={ids}
          songId={'song-1'}
          opened={true}
          onClose={onClose}
          onDelete={onDelete}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: /yes/i }))

    expect(capturedRequest).toStrictEqual({ ids: ['1', '2'], songId: 'song-1' })
    expect(screen.getByText('2 parts deleted!')).toBeInTheDocument()
    expect(onClose).toHaveBeenCalledOnce()
    expect(onDelete).toHaveBeenCalledOnce()
  })
})
