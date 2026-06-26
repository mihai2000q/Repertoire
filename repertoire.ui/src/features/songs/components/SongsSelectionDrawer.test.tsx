import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { useDragSelect } from '../../context/DragSelectContext.tsx'
import { emptySong, reduxRender } from '../../test-utils.tsx'
import SongsSelectionDrawer from './SongsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Song from '../../types/models/Song.ts'

// Mock the context
vi.mock('../../context/DragSelectContext', () => ({
  useDragSelect: vi.fn()
}))

describe('Songs Selection Drawer', () => {
  const selectedIds = ['1', '2', '3']
  const clearSelection = vi.fn()

  const songs: Song[] = [
    { ...emptySong, id: selectedIds[0] },
    { ...emptySong, id: selectedIds[1] },
    { ...emptySong, id: selectedIds[2] }
  ]

  const handlers = [
    http.get('/playlists', async () => {
      return HttpResponse.json([])
    }),
    http.get('/songs', async () => {
      return HttpResponse.json({ models: songs, totalCount: songs.length })
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
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<SongsSelectionDrawer />)

    expect(screen.getByText(`${selectedIds.length} songs selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
  })

  describe('on action icons', () => {
    it('should open warning when clicking on delete button', async () => {
      const user = userEvent.setup()

      reduxRender(<SongsSelectionDrawer />)

      await user.click(screen.getByRole('button', { name: 'delete' }))

      expect(await screen.findByRole('dialog', { name: /delete songs/i })).toBeInTheDocument()
    })
  })

  describe('on more menu', () => {
    it('should open custom rehearsals when clicking on custom rehearsals', async () => {
      const user = userEvent.setup()

      reduxRender(<SongsSelectionDrawer />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /custom rehearsals/i }))

      expect(await screen.findByRole('dialog', { name: /custom rehearsals/i })).toBeInTheDocument()
    })
  })
})
