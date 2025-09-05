import { reduxRender } from '../../test-utils.tsx'
import ArtistSongsContextMenu from './ArtistSongsContextMenu.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

describe('Artists Songs Context Menu', () => {
  const artistId = 'artist-id'
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
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
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
    vi.mock('../../context/ClickSelectContext', () => ({
      useClickSelect: vi.fn()
    }))
  })

  afterAll(() => server.close())

  const render = (isUnknownArtist: boolean = false) =>
    reduxRender(
      <ArtistSongsContextMenu artistId={artistId} isUnknownArtist={isUnknownArtist}>
        <div data-testid={dataTestId} />
      </ArtistSongsContextMenu>
    )

  it('should render with artist (not unknown)', async () => {
    const user = userEvent.setup()

    render()

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should render with unknown artist', async () => {
    const user = userEvent.setup()

    render(true)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByTestId(dataTestId)
    })

    expect(await screen.findByRole('menu')).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()

    expect(screen.queryByRole('menuitem', { name: /remove from artist/i })).not.toBeInTheDocument()
  })

  it('should be disabled when there are no selected ids', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
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
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      clearSelection: clearSelection
    })

    rerender(
      <ArtistSongsContextMenu artistId={artistId} isUnknownArtist={false}>
        <div data-testid={dataTestId} />
      </ArtistSongsContextMenu>
    )

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open warning when clicking on remove from artist menu item', async () => {
      const user = userEvent.setup()

      render()

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByTestId(dataTestId)
      })
      await user.click(screen.getByRole('menuitem', { name: /remove from artist/i }))

      expect(
        await screen.findByRole('dialog', { name: /remove songs from artist/i })
      ).toBeInTheDocument()
    })

    it('should open warning when clicking on delete menu item', async () => {
      const user = userEvent.setup()

      render()

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByTestId(dataTestId)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete songs/i })).toBeInTheDocument()
    })
  })
})
