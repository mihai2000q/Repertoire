import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import { reduxRender } from '../../test-utils.tsx'
import AlbumsSelectionDrawer from './AlbumsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

// Mock the context
vi.mock('../../context/DragSelectContext', () => ({
  useDragSelect: vi.fn()
}))

describe('Albums Selection Drawer', () => {
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

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<AlbumsSelectionDrawer />)

    expect(screen.getByText(`${selectedIds.length} albums selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
  })

  it('should open warning when clicking on delete button', async () => {
    const user = userEvent.setup()

    reduxRender(<AlbumsSelectionDrawer />)

    await user.click(screen.getByRole('button', { name: 'delete' }))

    expect(await screen.findByRole('dialog', { name: /delete albums/i })).toBeInTheDocument()
  })
})
