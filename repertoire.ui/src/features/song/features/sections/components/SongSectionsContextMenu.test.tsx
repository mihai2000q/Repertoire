import {
  emptySongPart,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../../../test-utils.tsx'
import SongSectionsContextMenu from './SongSectionsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { BulkUpdateSongPartsRequest } from '../../parts/types/requests/SongPartRequests.ts'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { SongSection } from '../../../../../types/models/Song.ts'

vi.mock('../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Sections Context Menu', () => {
  const dataTestId = 'dataTestId'
  const clearSelection = vi.fn()
  const sections: SongSection[] = ['1', '2', '3'].map((id) => ({
    ...emptySongSection,
    id,
    name: `Section ${id}`,
    parts: [{ ...emptySongPart, id: `part-${id}`, rehearsals: 2, confidence: 50 }]
  }))
  const server = setupServer()

  function mockSelectedIds(selectedIds: string[], isClickSelectionActive = true) {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds,
      isClickSelectionActive,
      clearSelection
    })
  }

  beforeEach(() => {
    vi.clearAllMocks()
    mockSelectedIds(['section-1', 'section-2', 'section-3'])
  })

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  const render = (songId = '1') =>
    reduxRender(
      withToastify(
        <SongSectionsContextMenu sections={sections} songId={songId}>
          <div data-testid={dataTestId} />
        </SongSectionsContextMenu>
      )
    )

  async function openMenu(songId = '1') {
    render(songId)
    const user = userEvent.setup()
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    return user
  }

  it('should open the menu for an active section selection', async () => {
    await openMenu()

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should not open when there is no active selection', async () => {
    mockSelectedIds([], false)

    await openMenu()

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should close when the selection becomes empty', async () => {
    const [{ rerender }] = render()
    const user = userEvent.setup()
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(screen.getByRole('menu')).toBeInTheDocument()

    mockSelectedIds([], false)
    rerender(
      <SongSectionsContextMenu sections={sections} songId={'1'}>
        <div data-testid={dataTestId} />
      </SongSectionsContextMenu>
    )
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should disable add rehearsals when only sections are selected', async () => {
    await openMenu()

    expect(screen.getByRole('menuitem', { name: /add rehearsals/i })).toHaveAttribute(
      'data-disabled'
    )
  })

  it('should open the section deletion modal when clicking delete', async () => {
    const user = await openMenu()

    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(await screen.findByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
  })

  it('should open the part deletion modal when clicking delete with only parts selected', async () => {
    const selectedPartIds = sections
      .slice(0, 2)
      .map((section) => `part-${section.parts[0].id}:${section.id}`)
    mockSelectedIds(selectedPartIds)

    const user = await openMenu()

    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(await screen.findByRole('dialog', { name: /delete parts/i })).toBeInTheDocument()
  })

  it('should open the sections and parts deletion modal when clicking delete with both selected', async () => {
    const selectedPartIds = sections
      .slice(0, 2)
      .map((section) => `part-${section.parts[0].id}:${section.id}`)
    mockSelectedIds([
      ...sections.slice(0, 2).map((section) => `section-${section.id}`),
      ...selectedPartIds
    ])

    const user = await openMenu()

    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(
      await screen.findByRole('dialog', { name: /delete song sections with parts/i })
    ).toBeInTheDocument()
    expect(screen.queryByRole('dialog', { name: /delete parts/i })).not.toBeInTheDocument()
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

    const songId = 'song-1'

    render(songId)
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    await user.click(screen.getByRole('menuitem', { name: /add rehearsals/i }))
    await user.click(screen.getByRole('button', { name: /confirm/i }))

    expect(
      screen.getByText(`Rehearsals added to ${selectedPartIds.length} parts!`)
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
