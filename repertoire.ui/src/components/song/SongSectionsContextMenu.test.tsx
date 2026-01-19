import { reduxRender, withToastify } from '../../test-utils.tsx'
import SongSectionsContextMenu from './SongSectionsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { BulkRehearsalsSongSectionsRequest } from '../../types/requests/SongRequests.ts'

describe('Song Sections Context Menu', () => {
  const dataTestId = 'dataTestId'
  const selectedIds = ['1', '2', '3']
  const clearSelection = vi.fn()

  const handlers = [
    http.get('/playlists', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: selectedIds,
      isSelectionActive: true,
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
    vi.mock('../../context/DragSelectContext', () => ({
      useDragSelect: vi.fn()
    }))
  })

  afterAll(() => server.close())

  const render = (songId = '1') =>
    reduxRender(
      withToastify(
        <SongSectionsContextMenu songId={songId}>
          <div data-testid={dataTestId} />
        </SongSectionsContextMenu>
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

  it('should be disabled when there are no selected ids', async () => {
    const user = userEvent.setup()

    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: [],
      clearSelection: clearSelection
    })

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should close menu when selected ids deplete', async () => {
    const user = userEvent.setup()

    // render and open menu
    const [{ rerender }] = render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    expect(screen.queryByRole('menu')).toBeInTheDocument()

    // empty the selected ids and rerender the closed menu
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: [],
      clearSelection: clearSelection
    })

    rerender(
      <SongSectionsContextMenu songId={'1'}>
        <div data-testid={dataTestId} />
      </SongSectionsContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should bulk rehearsals by 1 on add rehearsals menu item', async () => {
    const user = userEvent.setup()

    let capturedRequest: BulkRehearsalsSongSectionsRequest
    server.use(
      http.post(`/songs/sections/bulk-rehearsals`, async (req) => {
        capturedRequest = (await req.request.json()) as BulkRehearsalsSongSectionsRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = '1'

    render(songId)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    await user.click(screen.getByRole('menuitem', { name: /add rehearsals/i }))
    await user.click(screen.getByRole('button', { name: /confirm/i })) // menu item confirmation

    expect(
      screen.getByText(`Rehearsals added to ${selectedIds.length} sections!`)
    ).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({
      sections: selectedIds.map((id) => ({ id: id, rehearsals: 1 })),
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

    expect(await screen.findByRole('dialog', { name: /delete sections/i })).toBeInTheDocument()
  })
})
