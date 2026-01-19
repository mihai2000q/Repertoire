import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import { emptySong, reduxRender } from '../../test-utils.tsx'
import PlaylistSongsSelectionDrawer from './PlaylistSongsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Song from '../../types/models/Song.ts'

describe('Playlist Songs Selection Drawer', () => {
  const playlistId = 'playlist-id'
  const selectedIds = ['1', '2', '3']
  const clearSelection = vi.fn()

  const songs: Song[] = [
    { ...emptySong, id: 's1', playlistSongId: selectedIds[0] },
    { ...emptySong, id: 's2', playlistSongId: selectedIds[1] },
    { ...emptySong, id: 's3', playlistSongId: selectedIds[2] }
  ]

  const handlers = [
    http.get('/playlists', async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
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
    vi.mock('../../context/ClickSelectContext', () => ({
      useClickSelect: vi.fn()
    }))
  })

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistSongsSelectionDrawer playlistId={playlistId} songs={songs} />)

    expect(screen.getByText(`${selectedIds.length} songs selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'remove-from-playlist' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
  })

  describe('on action icons', () => {
    it('should open warning when clicking on remove from playlist button', async () => {
      const user = userEvent.setup()

      reduxRender(<PlaylistSongsSelectionDrawer playlistId={playlistId} songs={songs} />)

      await user.click(screen.getByRole('button', { name: 'remove-from-playlist' }))

      expect(
        await screen.findByRole('dialog', { name: /remove songs from playlist/i })
      ).toBeInTheDocument()
    })
  })
})
