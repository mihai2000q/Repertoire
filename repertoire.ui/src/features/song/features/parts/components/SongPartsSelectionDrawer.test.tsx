import { setupServer } from 'msw/node'
import { emptySongPart, reduxRender, withToastify } from '../../../../../test-utils.tsx'
import SongPartsSelectionDrawer from './SongPartsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { BulkUpdateSongPartsRequest } from '../types/requests/SongPartRequests.ts'
import { http, HttpResponse } from 'msw'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { SongPart } from '../../../../../types/models/Song.ts'

// Mock the context
vi.mock('../../../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Parts Selection Drawer', () => {
  const selectedIds = ['1', '3']
  const clearSelection = vi.fn()

  const parts: SongPart[] = [
    { ...emptySongPart, id: '1', rehearsals: 10, confidence: 2 },
    { ...emptySongPart, id: '2', rehearsals: 5, confidence: 10 },
    { ...emptySongPart, id: '3', rehearsals: 20, confidence: 50 }
  ]

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
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<SongPartsSelectionDrawer parts={parts} songId={'1'} />)

    expect(screen.getByText(`${selectedIds.length} parts selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsals' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
  })

  it('should bulk rehearsals by 1 when clicking on add rehearsals button', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkUpdateSongPartsRequest
    server.use(
      http.put(`/songs/parts/bulk-update`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkUpdateSongPartsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = '1'
    const selectedParts = parts.filter((p) => selectedIds.some((pId) => pId === p.id))

    reduxRender(withToastify(<SongPartsSelectionDrawer parts={parts} songId={songId} />))

    await user.click(screen.getByRole('button', { name: 'add-rehearsals' }))

    expect(screen.getByText(`Rehearsals added to ${selectedIds.length} parts!`)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      requests: selectedParts.map((p) => ({
        id: p.id,
        rehearsals: p.rehearsals + 1,
        confidence: p.confidence
      })),
      songId: songId
    })
    expect(clearSelection).toHaveBeenCalledOnce()
  })

  it('should open warning when clicking on delete button', async () => {
    const user = userEvent.setup()

    reduxRender(<SongPartsSelectionDrawer parts={parts} songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete parts/i })).toBeInTheDocument()
  })
})
