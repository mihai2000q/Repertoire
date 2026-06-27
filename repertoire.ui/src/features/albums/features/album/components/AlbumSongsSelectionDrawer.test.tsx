import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { emptySong, reduxRender } from '../../../../../test-utils.tsx'
import AlbumSongsSelectionDrawer from './AlbumSongsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Song from '../../../../../types/models/Song.ts'

// Mock the context
vi.mock('../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Album Songs Selection Drawer', () => {
  const albumId = 'album-id'
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
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: selectedIds,
      isClickSelectionActive: false,
      clearSelection: clearSelection
    })
  })

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  it('should render with album (not unknown)', async () => {
    const user = userEvent.setup()

    reduxRender(<AlbumSongsSelectionDrawer albumId={albumId} isUnknownAlbum={false} />)

    expect(screen.getByText(`${selectedIds.length} songs selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'remove-from-album' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /custom rehearsals/i })).toBeInTheDocument()
  })

  it('should render with unknown album', async () => {
    const user = userEvent.setup()

    reduxRender(<AlbumSongsSelectionDrawer albumId={albumId} isUnknownAlbum={true} />)

    expect(screen.getByText(`${selectedIds.length} songs selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    // on menu
    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /custom rehearsals/i })).toBeInTheDocument()

    // not available
    expect(screen.queryByRole('button', { name: 'remove-from-album' })).not.toBeInTheDocument()
  })

  describe('on action icons', () => {
    it('should open warning when clicking on remove from album button', async () => {
      const user = userEvent.setup()

      reduxRender(<AlbumSongsSelectionDrawer albumId={albumId} isUnknownAlbum={false} />)

      await user.click(screen.getByRole('button', { name: 'remove-from-album' }))

      expect(
        await screen.findByRole('dialog', { name: /remove songs from album/i })
      ).toBeInTheDocument()
    })

    it('should open warning when clicking on delete button', async () => {
      const user = userEvent.setup()

      reduxRender(<AlbumSongsSelectionDrawer albumId={albumId} isUnknownAlbum={false} />)

      await user.click(screen.getByRole('button', { name: 'delete' }))

      expect(await screen.findByRole('dialog', { name: /delete songs/i })).toBeInTheDocument()
    })
  })

  describe('on more menu', () => {
    it('should open custom rehearsals when clicking on custom rehearsals', async () => {
      const user = userEvent.setup()

      reduxRender(<AlbumSongsSelectionDrawer albumId={albumId} isUnknownAlbum={false} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /custom rehearsals/i }))

      expect(await screen.findByRole('dialog', { name: /custom rehearsals/i })).toBeInTheDocument()
    })
  })
})
