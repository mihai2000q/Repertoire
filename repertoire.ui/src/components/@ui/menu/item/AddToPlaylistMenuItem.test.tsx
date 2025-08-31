import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../../../types/models/Playlist.ts'
import { setupServer } from 'msw/node'
import { emptyPlaylist, reduxRender, withToastify } from '../../../../test-utils.tsx'
import { Menu } from '@mantine/core'
import AddToPlaylistMenuItem from './AddToPlaylistMenuItem.tsx'
import { fireEvent, screen, waitFor } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import plural from '../../../../utils/plural.ts'
import { describe, expect } from 'vitest'
import {
  AddAlbumsToPlaylistRequest,
  AddArtistsToPlaylistRequest,
  AddSongsToPlaylistRequest
} from '../../../../types/requests/PlaylistRequests.ts'
import {
  AddAlbumsToPlaylistResponse,
  AddArtistsToPlaylistResponse,
  AddSongsToPlaylistResponse
} from '../../../../types/responses/PlaylistResponses.ts'

describe('Add To Playlist Menu Item', () => {
  const playlists: Playlist[] = [
    {
      ...emptyPlaylist,
      id: '1',
      title: 'Playlist 1',
      songsCount: 12
    },
    {
      ...emptyPlaylist,
      id: '2',
      title: 'Playlist 2',
      imageUrl: 'something.png',
      songsCount: 1
    }
  ]

  const menuTargetId = 'target-id'
  const render = (props: {
    ids: string[]
    type: 'songs' | 'albums' | 'artists'
    closeMenu: () => void
    disabled?: boolean
  }) =>
    reduxRender(
      withToastify(
        <Menu>
          <Menu.Target>
            <button data-testid={menuTargetId}>Button</button>
          </Menu.Target>

          <Menu.Dropdown>
            <AddToPlaylistMenuItem {...props} />
          </Menu.Dropdown>
        </Menu>
      )
    )

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = {
        models: playlists,
        totalCount: playlists.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    render({ ids: [], type: 'songs', closeMenu: vi.fn() })

    await user.click(screen.getByTestId(menuTargetId))
    const menuItem = screen.getByRole('menuitem', { name: /add to playlist/i })
    expect(menuItem).toBeInTheDocument()

    await user.hover(menuItem)
    expect(screen.getByRole('textbox', { name: /search/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /search/i })).toHaveValue('')

    for (const playlist of playlists) {
      expect(screen.getByRole('menuitem', { name: playlist.title })).toBeInTheDocument()

      if (playlist.imageUrl)
        expect(screen.getByRole('img', { name: playlist.title })).toBeInTheDocument()
      else expect(screen.getByLabelText(`default-icon-${playlist.title}`)).toBeInTheDocument()

      expect(screen.getByText(playlist.title)).toBeInTheDocument()
      expect(
        screen.getByText(`${playlist.songsCount} song${plural(playlist.songsCount)}`)
      ).toBeInTheDocument()
    }
  })

  it('should not be disabled when there are playlists', async () => {
    const user = userEvent.setup()

    render({ ids: [], type: 'songs', closeMenu: vi.fn() })

    await user.click(screen.getByTestId(menuTargetId))
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).not.toBeDisabled()
  })

  describe('should be disabled', () => {
    it('when there are no playlists', async () => {
      const user = userEvent.setup()

      server.use(
        http.get('/playlists', async () => {
          const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
          return HttpResponse.json(response)
        })
      )

      render({ ids: [], type: 'songs', closeMenu: vi.fn() })

      await user.click(screen.getByTestId(menuTargetId))
      expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
      expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeDisabled()
    })

    it('when the props contain disabled', async () => {
      const user = userEvent.setup()

      render({ ids: [], type: 'songs', closeMenu: vi.fn(), disabled: true })

      await user.click(screen.getByTestId(menuTargetId))
      expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
      expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeDisabled()
    })
  })

  describe('should send add request when selecting a playlist with success', () => {
    it('when type is song', async () => {
      const user = userEvent.setup()

      const ids = ['id-1']
      const closeMenu = vi.fn()
      const newPlaylist = playlists[1]

      let capturedRequest: AddSongsToPlaylistRequest
      server.use(
        http.post('/playlists/songs/add', async (req) => {
          capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
          const response: AddSongsToPlaylistResponse = {
            success: true,
            duplicates: [],
            added: ids
          }
          return HttpResponse.json(response)
        })
      )

      render({ ids: ids, type: 'songs', closeMenu: closeMenu })

      await user.click(screen.getByTestId(menuTargetId))
      await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
      fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

      expect(
        await screen.findByText(new RegExp(`successfully added ${ids.length} song`, 'i'))
      ).toBeInTheDocument()
      expect(closeMenu).toHaveBeenCalledOnce()
      expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
    })

    it('when type is album', async () => {
      const user = userEvent.setup()

      const ids = ['id-1']
      const closeMenu = vi.fn()
      const newPlaylist = playlists[1]

      const addedSongIds = ['songId1', 'songId2']

      let capturedRequest: AddAlbumsToPlaylistRequest
      server.use(
        http.post('/playlists/add-albums', async (req) => {
          capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
          const response: AddAlbumsToPlaylistResponse = {
            success: true,
            duplicateAlbumIds: [],
            duplicateSongIds: [],
            addedSongIds: addedSongIds
          }
          return HttpResponse.json(response)
        })
      )

      render({ ids: ids, type: 'albums', closeMenu: closeMenu })

      await user.click(screen.getByTestId(menuTargetId))
      await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
      fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

      expect(
        await screen.findByText(new RegExp(`successfully added ${addedSongIds.length} songs`, 'i'))
      ).toBeInTheDocument()
      expect(closeMenu).toHaveBeenCalledOnce()
      expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
    })

    it('when type is artist', async () => {
      const user = userEvent.setup()

      const ids = ['id-1']
      const closeMenu = vi.fn()
      const newPlaylist = playlists[1]

      const addedSongIds = ['songId1', 'songId2']

      let capturedRequest: AddArtistsToPlaylistRequest
      server.use(
        http.post('/playlists/add-artists', async (req) => {
          capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
          const response: AddArtistsToPlaylistResponse = {
            success: true,
            duplicateArtistIds: [],
            duplicateSongIds: [],
            addedSongIds: addedSongIds
          }
          return HttpResponse.json(response)
        })
      )

      render({ ids: ids, type: 'artists', closeMenu: closeMenu })

      await user.click(screen.getByTestId(menuTargetId))
      await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
      fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

      expect(
        await screen.findByText(new RegExp(`successfully added ${addedSongIds.length} songs`, 'i'))
      ).toBeInTheDocument()
      expect(closeMenu).toHaveBeenCalledOnce()
      expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
    })
  })

  describe('should display already added modal when selecting a playlist with failed success due to duplicates', () => {
    describe('when type is song', () => {
      it("when there is only one song and it's duplicated and it cancels", async () => {
        const user = userEvent.setup()

        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: ids,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/this song is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        await waitFor(() => expect(closeMenu).toHaveBeenCalledOnce())
      })

      it("when there is only one song and it's duplicated and it adds it anyway", async () => {
        const user = userEvent.setup()

        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddSongsToPlaylistResponse = {
                success: true,
                duplicates: ids,
                added: ids
              }
              return HttpResponse.json(response)
            }
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: ids,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/this song is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(new RegExp(`successfully added ${ids.length} song`, 'i'))
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs and it adds them all', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3']
        const duplicateIds = ['id-1', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddSongsToPlaylistResponse = {
                success: true,
                duplicates: duplicateIds,
                added: ids
              }
              return HttpResponse.json(response)
            }
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: duplicateIds,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add all/i }))

        expect(
          await screen.findByText(new RegExp(`successfully added ${ids.length} songs`, 'i'))
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs and it adds only the new ones', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const duplicateIds = ['id-1', 'id-3']
        const addedIds = ['id-2', 'id-4']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddSongsToPlaylistResponse = {
                success: true,
                duplicates: duplicateIds,
                added: addedIds
              }
              return HttpResponse.json(response)
            }
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: duplicateIds,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /just new ones/i }))

        expect(
          await screen.findByText(new RegExp(`successfully added ${addedIds.length} songs`, 'i'))
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids, forceAdd: false })
      })

      it('when these are all duplicated songs and it cancels', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: ids,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/these songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        expect(closeMenu).toHaveBeenCalledOnce()
      })

      it('when these are all duplicated songs and it adds them anyway', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddSongsToPlaylistRequest
        server.use(
          http.post('/playlists/songs/add', async (req) => {
            capturedRequest = (await req.request.json()) as AddSongsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddSongsToPlaylistResponse = {
                success: true,
                duplicates: ids,
                added: ids
              }
              return HttpResponse.json(response)
            }
            const response: AddSongsToPlaylistResponse = {
              success: false,
              duplicates: ids,
              added: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'songs', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids })
        expect(screen.getByText(/these songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(new RegExp(`successfully added ${ids.length} songs`, 'i'))
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, songIds: ids, forceAdd: true })
      })
    })

    describe('when type is album', () => {
      it("when there is only one album and it's duplicated and it cancels", async () => {
        const user = userEvent.setup()

        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/this album is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        await waitFor(() => expect(closeMenu).toHaveBeenCalledOnce())
      })

      it("when there is only one album and it's duplicated and it adds it anyway", async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: false,
                duplicateAlbumIds: ids,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/this album is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids, forceAdd: true })
      })

      it('when there are some duplicated albums and it adds them all', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const duplicateIds = ['id-1', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: true,
                duplicateAlbumIds: duplicateIds,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: duplicateIds,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/some albums are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add all/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids, forceAdd: true })
      })

      it('when there are some duplicated albums and it adds only the new ones', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const duplicateIds = ['id-1', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: true,
                duplicateAlbumIds: duplicateIds,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: duplicateIds,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/some albums are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /just new ones/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({
          id: newPlaylist.id,
          albumIds: ids,
          forceAdd: false
        })
      })

      it('when these are all duplicated albums and it cancels', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/these albums are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        expect(closeMenu).toHaveBeenCalledOnce()
      })

      it('when these are all duplicated albums and it adds them anyway', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: true,
                duplicateAlbumIds: ids,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/these albums are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs in the albums and it adds them all', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const duplicateIds = ['song-id-1', 'song-id-2']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: true,
                duplicateAlbumIds: [],
                duplicateSongIds: duplicateIds,
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: [],
              duplicateSongIds: duplicateIds,
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add all/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs in the albums and it adds only the new ones', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const duplicateIds = ['song-id-2', 'song-id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddAlbumsToPlaylistRequest
        server.use(
          http.post('/playlists/add-albums', async (req) => {
            capturedRequest = (await req.request.json()) as AddAlbumsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddAlbumsToPlaylistResponse = {
                success: true,
                duplicateAlbumIds: [],
                duplicateSongIds: duplicateIds,
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddAlbumsToPlaylistResponse = {
              success: false,
              duplicateAlbumIds: [],
              duplicateSongIds: duplicateIds,
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'albums', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, albumIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /just new ones/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({
          id: newPlaylist.id,
          albumIds: ids,
          forceAdd: false
        })
      })
    })

    describe('when type is artist', () => {
      it("when there is only one artist and it's duplicated and it cancels", async () => {
        const user = userEvent.setup()

        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/this artist is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        await waitFor(() => expect(closeMenu).toHaveBeenCalledOnce())
      })

      it("when there is only one artist and it's duplicated and it adds it anyway", async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: false,
                duplicateArtistIds: ids,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/this artist is already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids, forceAdd: true })
      })

      it('when there are some duplicated artists and it adds them all', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const duplicateIds = ['id-1', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: true,
                duplicateArtistIds: duplicateIds,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: duplicateIds,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/some artists are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add all/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids, forceAdd: true })
      })

      it('when there are some duplicated artists and it adds only the new ones', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const duplicateIds = ['id-1', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: true,
                duplicateArtistIds: duplicateIds,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: duplicateIds,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/some artists are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /just new ones/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({
          id: newPlaylist.id,
          artistIds: ids,
          forceAdd: false
        })
      })

      it('when these are all duplicated artists and it cancels', async () => {
        const user = userEvent.setup()

        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/these artists are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /cancel/i }))

        expect(closeMenu).toHaveBeenCalledOnce()
      })

      it('when these are all duplicated artists and it adds them anyway', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: true,
                duplicateArtistIds: ids,
                duplicateSongIds: [],
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: ids,
              duplicateSongIds: [],
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/these artists are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add anyway/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs in the artists and it adds them all', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-2', 'song-id-3', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3']
        const duplicateIds = ['song-id-1', 'song-id-2']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: true,
                duplicateArtistIds: [],
                duplicateSongIds: duplicateIds,
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: [],
              duplicateSongIds: duplicateIds,
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /add all/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids, forceAdd: true })
      })

      it('when there are some duplicated songs in the artists and it adds only the new ones', async () => {
        const user = userEvent.setup()

        const addedSongIds = ['song-id-1', 'song-id-4']
        const ids = ['id-1', 'id-2', 'id-3', 'id-4']
        const duplicateIds = ['song-id-2', 'song-id-3']
        const closeMenu = vi.fn()
        const newPlaylist = playlists[1]

        let capturedRequest: AddArtistsToPlaylistRequest
        server.use(
          http.post('/playlists/add-artists', async (req) => {
            capturedRequest = (await req.request.json()) as AddArtistsToPlaylistRequest
            if (capturedRequest.forceAdd !== undefined) {
              const response: AddArtistsToPlaylistResponse = {
                success: true,
                duplicateArtistIds: [],
                duplicateSongIds: duplicateIds,
                addedSongIds: addedSongIds
              }
              return HttpResponse.json(response)
            }
            const response: AddArtistsToPlaylistResponse = {
              success: false,
              duplicateArtistIds: [],
              duplicateSongIds: duplicateIds,
              addedSongIds: []
            }
            return HttpResponse.json(response)
          })
        )

        render({ ids: ids, type: 'artists', closeMenu: closeMenu })

        await user.click(screen.getByTestId(menuTargetId))
        await user.hover(screen.getByRole('menuitem', { name: /add to playlist/i }))
        fireEvent.click(screen.getByRole('menuitem', { name: newPlaylist.title }))

        expect(await screen.findByRole('dialog', { name: /already added/i })).toBeInTheDocument()
        expect(capturedRequest).toStrictEqual({ id: newPlaylist.id, artistIds: ids })
        expect(screen.getByText(/some songs are already/i)).toBeInTheDocument()

        await user.click(screen.getByRole('button', { name: /just new ones/i }))

        expect(
          await screen.findByText(
            new RegExp(`successfully added ${addedSongIds.length} songs`, 'i')
          )
        ).toBeInTheDocument()
        expect(closeMenu).toHaveBeenCalledOnce()
        expect(capturedRequest).toStrictEqual({
          id: newPlaylist.id,
          artistIds: ids,
          forceAdd: false
        })
      })
    })
  })
})
