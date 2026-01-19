import { setupServer } from 'msw/node'
import { reduxRender, withToastify } from '../../test-utils.tsx'
import SongSectionsSelectionDrawer from './SongSectionsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { BulkRehearsalsSongSectionsRequest } from '../../types/requests/SongRequests.ts'
import { http, HttpResponse } from 'msw'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

describe('Song Sections Selection Drawer', () => {
  const selectedIds = ['1', '2', '3']
  const clearSelection = vi.fn()

  const server = setupServer()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: selectedIds,
      isClickSelectionActive: true,
      clearSelection: clearSelection
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
    server.resetHandlers()
  })

  beforeAll(() => {
    server.listen()
    // Mock the context
    vi.mock('../../context/ClickSelectContext', () => ({
      useDragSelect: vi.fn()
    }))
  })

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<SongSectionsSelectionDrawer songId={'1'} />)

    expect(screen.getByText(`${selectedIds.length} sections selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsals' })).toBeInTheDocument()
  })

  it('should open warning when clicking on delete button', async () => {
    const user = userEvent.setup()

    reduxRender(<SongSectionsSelectionDrawer songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
  })

  it('should bulk rehearsals by 1 when clicking on add rehearsals button', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkRehearsalsSongSectionsRequest
    server.use(
      http.post(`/songs/sections/bulk-rehearsals`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkRehearsalsSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = '1'

    reduxRender(withToastify(<SongSectionsSelectionDrawer songId={songId} />))

    await user.click(screen.getByRole('button', { name: 'add-rehearsals' }))

    expect(
      screen.getByText(`Rehearsals added to ${selectedIds.length} sections!`)
    ).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      sections: selectedIds.map((id) => ({ id: id, rehearsals: 1 })),
      songId: songId
    })
    expect(clearSelection).toHaveBeenCalledOnce()
  })
})
