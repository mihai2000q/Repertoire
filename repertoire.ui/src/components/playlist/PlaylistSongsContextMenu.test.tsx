import { emptySong, reduxRender } from '../../test-utils.tsx'
import PlaylistSongsContextMenu from './PlaylistSongsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import Song from '../../types/models/Song.ts'

describe('Playlists Songs Context Menu', () => {
  const playlistId = 'playlist-id'
  const dataTestId = 'dataTestId'
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
      isClickSelectionActive: true,
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

  const render = () =>
    reduxRender(
      <PlaylistSongsContextMenu playlistId={playlistId} songs={songs}>
        <div data-testid={dataTestId} />
      </PlaylistSongsContextMenu>
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
    expect(screen.getByRole('menuitem', { name: /remove from playlist/i })).toBeInTheDocument()
  })

  it('should be disabled when the selection is inactive', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: false,
      clearSelection: clearSelection
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
      clearSelection: clearSelection
    })

    rerender(
      <PlaylistSongsContextMenu playlistId={playlistId} songs={songs}>
        <div data-testid={dataTestId} />
      </PlaylistSongsContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open warning when clicking on remove from playlist menu item', async () => {
      const user = userEvent.setup()

      render()

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByTestId(dataTestId)
      })
      await user.click(screen.getByRole('menuitem', { name: /remove from playlist/i }))

      expect(
        await screen.findByRole('dialog', { name: /remove songs from playlist/i })
      ).toBeInTheDocument()
    })
  })
})
