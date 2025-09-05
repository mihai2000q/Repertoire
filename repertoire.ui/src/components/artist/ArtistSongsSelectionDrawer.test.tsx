import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'
import { reduxRender } from '../../test-utils.tsx'
import ArtistAlbumsSelectionDrawer from './ArtistAlbumsSelectionDrawer.tsx'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'

describe('Artist Albums Selection Drawer', () => {
  const artistId = 'artist-id'
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

  it('should render with artist (not unknown)', async () => {
    const user = userEvent.setup()

    reduxRender(<ArtistAlbumsSelectionDrawer artistId={artistId} isUnknownArtist={false} />)

    expect(screen.getByText(`${selectedIds.length} albums selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'remove-from-artist' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()
  })

  it('should render with unknown artist', async () => {
    const user = userEvent.setup()

    reduxRender(<ArtistAlbumsSelectionDrawer artistId={artistId} isUnknownArtist={true} />)

    expect(screen.getByText(`${selectedIds.length} albums selected`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'delete' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    // on menu
    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(await screen.findByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsals/i })).toBeInTheDocument()

    // not available
    expect(screen.queryByRole('button', { name: 'remove-from-artist' })).not.toBeInTheDocument()
  })

  describe('on action icons', () => {
    it('should open warning when clicking on remove from artist button', async () => {
      const user = userEvent.setup()

      reduxRender(<ArtistAlbumsSelectionDrawer artistId={artistId} isUnknownArtist={false} />)

      await user.click(screen.getByRole('button', { name: 'remove-from-artist' }))

      expect(
        await screen.findByRole('dialog', { name: /remove albums from artist/i })
      ).toBeInTheDocument()
    })

    it('should open warning when clicking on delete button', async () => {
      const user = userEvent.setup()

      reduxRender(<ArtistAlbumsSelectionDrawer artistId={artistId} isUnknownArtist={false} />)

      await user.click(screen.getByRole('button', { name: 'delete' }))

      expect(await screen.findByRole('dialog', { name: /delete albums/i })).toBeInTheDocument()
    })
  })
})
