import { reduxRender } from '../../test-utils.tsx'
import PlaylistsContextMenu from './PlaylistsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { useDragSelect } from '../../context/DragSelectContext.tsx'

describe('Playlists Context Menu', () => {
  const dataTestId = 'dataTestId'
  const selectedIds = ['1', '2', '3']
  const clearSelection = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: selectedIds,
      clearSelection: clearSelection
    })
  })

  beforeAll(() => {
    // Mock the context
    vi.mock('../../context/DragSelectContext', () => ({
      useDragSelect: vi.fn()
    }))
  })

  afterEach(() => vi.restoreAllMocks())

  const render = () =>
    reduxRender(
      <PlaylistsContextMenu>
        <div data-testid={dataTestId} />
      </PlaylistsContextMenu>
    )

  it('should render', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
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
      <PlaylistsContextMenu>
        <div data-testid={dataTestId} />
      </PlaylistsContextMenu>
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

    expect(await screen.findByRole('dialog', { name: /delete playlists/i })).toBeInTheDocument()
  })
})
