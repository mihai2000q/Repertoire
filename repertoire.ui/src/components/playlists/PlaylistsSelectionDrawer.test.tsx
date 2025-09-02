import { useDragSelect } from '../../context/DragSelectContext.tsx'
import { reduxRender } from '../../test-utils.tsx'
import PlaylistsSelectionDrawer from './PlaylistsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Playlists Selection Drawer', () => {
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
    // Mock Context
    vi.mock('../../context/DragSelectContext', () => ({
      useDragSelect: vi.fn()
    }))
  })

  afterEach(() => vi.restoreAllMocks())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistsSelectionDrawer />)

    expect(screen.getByText(`${selectedIds.length} playlists selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
  })

  it('should open warning when clicking on delete button', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistsSelectionDrawer />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete playlists/i })).toBeInTheDocument()
  })
})
