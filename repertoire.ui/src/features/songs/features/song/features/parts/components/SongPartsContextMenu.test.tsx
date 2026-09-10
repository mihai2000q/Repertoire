import { emptySongPart, reduxRender, withToastify } from '../../../../../../../test-utils.tsx'
import SongPartsContextMenu from './SongPartsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { BulkUpdateSongPartsRequest } from '../types/requests/SongPartRequests.ts'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongPart } from '../../../../../../../types/models/Song.ts'

// Mock the context
vi.mock('../../../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Parts Context Menu', () => {
  const dataTestId = 'dataTestId'
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

  const render = (songId = '1') =>
    reduxRender(
      withToastify(
        <SongPartsContextMenu parts={parts} songId={songId}>
          <div data-testid={dataTestId} />
        </SongPartsContextMenu>
      )
    )

  it('should render', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should be disabled when the selection is inactive', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: false,
      clearSelection: vi.fn()
    })

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should close menu when the selection becomes inactive', async () => {
    const user = userEvent.setup()

    // render and open menu
    const [{ rerender }] = render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    expect(screen.queryByRole('menu')).toBeInTheDocument()

    // close the activity of the selection and rerender the closed menu
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: false,
      clearSelection: vi.fn()
    })

    rerender(
      <SongPartsContextMenu parts={parts} songId={'1'}>
        <div data-testid={dataTestId} />
      </SongPartsContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should bulk rehearsals by 1 on add rehearsals menu item', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkUpdateSongPartsRequest
    server.use(
      http.post(`/songs/parts/bulk-update`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkUpdateSongPartsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = '1'
    const selectedParts = parts.filter((p) => selectedIds.some((pId) => pId === p.id))

    render(songId)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    await user.click(screen.getByRole('menuitem', { name: /add rehearsals/i }))
    await user.click(screen.getByRole('button', { name: /confirm/i })) // menu item confirmation

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

  it('should open warning when clicking on delete menu item', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(await screen.findByRole('dialog', { name: /delete parts/i })).toBeInTheDocument()
  })
})
