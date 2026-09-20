import { setupServer } from 'msw/node'
import {
  emptySongPart,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../../../test-utils.tsx'
import SongSectionsSelectionDrawer from './SongSectionsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { BulkUpdateSongPartsRequest } from '../../parts/types/requests/SongPartRequests.ts'
import { http, HttpResponse } from 'msw'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { SongSection } from '../../../../../types/models/Song.ts'

// Mock the context
vi.mock('../../../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Sections Selection Drawer', () => {
  const selectedSectionIds = ['section-1', 'section-2', 'section-3']
  const clearSelection = vi.fn()
  const sections: SongSection[] = ['1', '2', '3'].map((id) => ({
    ...emptySongSection,
    id,
    name: `Section ${id}`,
    parts: [{ ...emptySongPart, id: `part-${id}`, rehearsals: 2, confidence: 50 }]
  }))

  const server = setupServer()

  function mockSelectedIds(selectedIds: string[]) {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: selectedIds,
      isClickSelectionActive: true,
      clearSelection: clearSelection
    })
  }

  beforeEach(() => {
    vi.clearAllMocks()
    mockSelectedIds(selectedSectionIds)
  })

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render the number of selected sections', () => {
    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    expect(screen.getByText(`${selectedSectionIds.length} sections selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsals' })).toBeInTheDocument()
  })

  it('should disable add rehearsals when only sections are selected', () => {
    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    expect(screen.getByRole('button', { name: 'add-rehearsals' })).toBeDisabled()
  })

  it('should clear the selection when closing the drawer', async () => {
    const user = userEvent.setup()

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'close-drawer' }))

    expect(clearSelection).toHaveBeenCalledOnce()
  })

  it('should not render when click selection is inactive', () => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: selectedSectionIds,
      isClickSelectionActive: false,
      clearSelection
    })

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    expect(screen.queryByText('3 sections selected')).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'delete' })).not.toBeInTheDocument()
  })

  it('should open the section deletion modal when clicking delete with sections selected', async () => {
    const user = userEvent.setup()

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
  })

  it('should open the part deletion modal when clicking delete with only parts selected', async () => {
    const user = userEvent.setup()
    const selectedPartIds = sections
      .slice(0, 2)
      .map((section) => `part-${section.parts[0].id}:${section.id}`)
    mockSelectedIds(selectedPartIds)

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete parts/i })).toBeInTheDocument()
  })

  it('should open the sections and parts deletion modal when clicking delete with both selected', async () => {
    const user = userEvent.setup()
    const selectedPartIds = sections
      .slice(0, 2)
      .map((section) => `part-${section.parts[0].id}:${section.id}`)
    mockSelectedIds([...selectedSectionIds.slice(0, 2), ...selectedPartIds])

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(
      await screen.findByRole('dialog', { name: /delete song sections with parts/i })
    ).toBeInTheDocument()
    expect(screen.queryByRole('dialog', { name: /delete parts/i })).not.toBeInTheDocument()
  })

  it('should display both selected sections and parts', () => {
    const selectedPartIds = sections
      .slice(0, 2)
      .map((section) => `part-${section.parts[0].id}:${section.id}`)
    mockSelectedIds([...selectedSectionIds.slice(0, 2), ...selectedPartIds])

    reduxRender(<SongSectionsSelectionDrawer sections={sections} songId={'1'} />)

    expect(
      screen.getByText(`${sections} sections and ${selectedPartIds} parts selected`)
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsals' })).toBeEnabled()
  })

  it("should bulk rehearsals' selected parts by 1", async () => {
    const user = userEvent.setup()

    const selectedPartIds = sections.map((section) => `part-${section.parts[0].id}:${section.id}`)
    const selectedParts = sections.map((section) => section.parts[0])
    mockSelectedIds(selectedPartIds)

    let capturedRequest: BulkUpdateSongPartsRequest
    server.use(
      http.put(`/songs/parts/bulk-update`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkUpdateSongPartsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = '1'

    reduxRender(withToastify(<SongSectionsSelectionDrawer sections={sections} songId={songId} />))

    await user.click(screen.getByRole('button', { name: 'add-rehearsals' }))

    expect(
      screen.getByText(`Rehearsals added to ${selectedParts.length} parts!`)
    ).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      requests: selectedParts.map((part) => ({
        id: part.id,
        rehearsals: part.rehearsals + 1,
        confidence: part.confidence
      })),
      songId: songId
    })
    expect(clearSelection).toHaveBeenCalledOnce()
  })
})
