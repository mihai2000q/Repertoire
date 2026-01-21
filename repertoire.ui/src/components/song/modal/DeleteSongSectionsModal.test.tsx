import { http, HttpResponse } from 'msw'
import { reduxRender, withToastify } from '../../../test-utils.tsx'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import DeleteSongSectionsModal from './DeleteSongSectionsModal.tsx'
import { BulkDeleteSongSectionsRequest } from '../../../types/requests/SongRequests.ts'

describe('Delete Song Sections Modal', () => {
  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(
      <DeleteSongSectionsModal ids={['1', '2']} songId={'1'} opened={true} onClose={vi.fn()} />
    )

    expect(screen.getByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /delete sections/i })).toBeInTheDocument()
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /yes/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('should delete song sections in bulk call onDelete', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkDeleteSongSectionsRequest
    server.use(
      http.put(`/songs/sections/bulk-delete`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkDeleteSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const ids = ['1', '2']
    const songId = '1'
    const onClose = vi.fn()
    const onDelete = vi.fn()

    reduxRender(
      withToastify(
        <DeleteSongSectionsModal
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
    expect(screen.getByText(`${ids.length} sections deleted!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ ids: ids, songId: songId })
  })
})
