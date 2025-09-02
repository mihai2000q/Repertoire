import Playlist from 'src/types/models/Playlist.ts'
import { emptyPlaylist, reduxRouterRender, withToastify } from '../../test-utils.tsx'
import PlaylistCard from './PlaylistCard.tsx'
import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { expect } from 'vitest'
import { RootState } from '../../state/store.ts'
import { useDragSelect } from '../../context/DragSelectContext.tsx'

describe('Playlist Card', () => {
  const playlist: Playlist = {
    ...emptyPlaylist,
    id: '1',
    title: 'Playlist 1'
  }

  const server = setupServer()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: [],
      clearSelection: vi.fn()
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

  it('should render minimal info', () => {
    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toBeInTheDocument()
    expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toHaveAttribute(
      'aria-selected',
      'false'
    )
    expect(screen.getByLabelText(`default-icon-${playlist.title}`)).toBeInTheDocument()
    expect(screen.getByText(playlist.title)).toBeInTheDocument()
  })

  it('should render maximal info', () => {
    const localPlaylist: Playlist = {
      ...playlist,
      imageUrl: 'something.png'
    }

    reduxRouterRender(<PlaylistCard playlist={localPlaylist} />)

    expect(screen.getByLabelText(`playlist-card-${localPlaylist.title}`)).toHaveAttribute(
      'aria-selected',
      'false'
    )
    expect(screen.getByRole('img', { name: localPlaylist.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localPlaylist.title })).toHaveAttribute(
      'src',
      localPlaylist.imageUrl
    )
    expect(screen.getByText(localPlaylist.title)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${playlist.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /open drawer/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open playlist drawer when clicking on open drawer', async () => {
      const user = userEvent.setup()

      const [_, store] = reduxRouterRender(<PlaylistCard playlist={playlist} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${playlist.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /open drawer/i }))

      expect((store.getState() as RootState).global.playlistDrawer.open).toBeTruthy()
      expect((store.getState() as RootState).global.playlistDrawer.playlistId).toBe(playlist.id)
    })

    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/playlists/${playlist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<PlaylistCard playlist={playlist} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${playlist.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${playlist.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    await user.click(screen.getByLabelText(`default-icon-${playlist.title}`))
    expect(window.location.pathname).toBe(`/playlist/${playlist.id}`)
  })

  it('should disable context menu there are selected ids', async () => {
    const user = userEvent.setup()

    vi.mocked(useDragSelect).mockReturnValue({
      dragSelect: null,
      selectedIds: ['someone'],
      clearSelection: vi.fn()
    })

    reduxRouterRender(<PlaylistCard playlist={playlist} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${playlist.title}`)
    })

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  describe('should be selected', () => {
    it('when avatar is hovered', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<PlaylistCard playlist={playlist} />)

      await user.hover(screen.getByLabelText(`default-icon-${playlist.title}`))

      expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<PlaylistCard playlist={playlist} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${playlist.title}`)
      })

      expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when drag selected (part of the selected ids)', () => {
      vi.mocked(useDragSelect).mockReturnValue({
        dragSelect: null,
        selectedIds: [playlist.id],
        clearSelection: vi.fn()
      })

      reduxRouterRender(<PlaylistCard playlist={playlist} />)

      expect(screen.getByLabelText(`playlist-card-${playlist.title}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
