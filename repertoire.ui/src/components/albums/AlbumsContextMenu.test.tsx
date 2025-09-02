import { reduxRender } from '../../test-utils.tsx'
import AlbumsContextMenu from './AlbumsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

describe('Albums Context Menu', () => {
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

  const render = () =>
    reduxRender(
      <AlbumsContextMenu>
        <div data-testid={dataTestId} />
      </AlbumsContextMenu>
    )

  it('should render', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
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
      <AlbumsContextMenu>
        <div data-testid={dataTestId} />
      </AlbumsContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should open warning when clicking on delete menu item', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })
    await user.click(screen.getByRole('menuitem', { name: /delete/i }))

    expect(await screen.findByRole('dialog', { name: /delete albums/i })).toBeInTheDocument()
  })
})
