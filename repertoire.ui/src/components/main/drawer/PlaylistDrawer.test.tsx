import {
  emptyAlbum,
  emptyPlaylist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../../test-utils.tsx'
import PlaylistDrawer from './PlaylistDrawer.tsx'
import Playlist from '../../../types/models/Playlist.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { act, screen, waitFor } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import { expect } from 'vitest'
import { openPlaylistDrawer, setDocumentTitle } from '../../../state/slice/globalSlice.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'

describe('Playlist Drawer', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      playlistSongId: '1',
      title: 'Song 1',
      playlistTrackNo: 1,
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      id: '1',
      playlistSongId: '2',
      title: 'Song 2',
      playlistTrackNo: 2,
      album: {
        ...emptyAlbum,
        imageUrl: 'something.png'
      }
    },
    {
      ...emptySong,
      id: '3',
      playlistSongId: '3',
      title: 'Song 3',
      playlistTrackNo: 3
    }
  ]

  const playlist: Playlist = {
    ...emptyPlaylist,
    id: '1',
    title: 'Playlist 1',
    songs: songs
  }

  const getPlaylist = (playlist: Playlist) =>
    http.get(`/playlists/${playlist.id}`, async () => {
      return HttpResponse.json(playlist)
    })

  const handlers = [
    getPlaylist(playlist),
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  const prevDocumentTitle = 'previous document title'

  const render = (id: string | null = playlist.id) =>
    reduxRouterRender(<PlaylistDrawer />, {
      global: {
        documentTitle: prevDocumentTitle,
        playlistDrawer: {
          open: true,
          playlistId: id
        },
        albumDrawer: undefined,
        artistDrawer: undefined,
        songDrawer: undefined
      }
    })

  it('should render and display minimal info', async () => {
    const [_, store] = render()

    expect(screen.getByTestId('playlist-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${playlist.title}`)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: playlist.title })).toBeInTheDocument()
    expect(screen.getByText(`${playlist.songs.length} songs`)).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + playlist.title
    )

    playlist.songs.forEach((song) => {
      expect(screen.getByText(song.playlistTrackNo)).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else if (song.album?.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      }
    })
  })

  it('should render and display maximal info', async () => {
    const localPlaylist = {
      ...playlist,
      imageUrl: 'something.png'
    }

    server.use(getPlaylist(localPlaylist))

    const [_, store] = render()

    expect(screen.getByTestId('playlist-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localPlaylist.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localPlaylist.title })).toHaveAttribute(
      'src',
      localPlaylist.imageUrl
    )
    expect(screen.getByRole('heading', { name: localPlaylist.title })).toBeInTheDocument()
    expect(screen.getByText(`${localPlaylist.songs.length} songs`)).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + localPlaylist.title
    )

    localPlaylist.songs.forEach((song) => {
      expect(screen.getByText(song.playlistTrackNo)).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else if (song.album?.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      }
    })
  })

  it('should render the loader when there is no playlist', async () => {
    render(null)

    expect(screen.getByTestId('playlist-drawer-loader')).toBeInTheDocument()
  })

  it('should change document title on opening and closing, correctly', async () => {
    const newPlaylist: Playlist = {
      ...playlist,
      id: 'another-id-another-playlist',
      title: 'new Title'
    }

    const [_, store] = render()

    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + playlist.title
      )
    })

    // change back the document title (as if the drawer closed)
    await act(() => store.dispatch(setDocumentTitle(prevDocumentTitle)))

    // make sure it doesn't use the old title when a new playlist is introduced
    server.use(getPlaylist(newPlaylist))
    await act(() => store.dispatch(openPlaylistDrawer(newPlaylist.id)))
    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + newPlaylist.title
      )
    })
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to playlist when clicking on view details', async () => {
      const user = userEvent.setup()

      const [_, store] = render()

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/playlist/${playlist.id}`)
      expect((store.getState() as RootState).global.playlistDrawer.open).toBeFalsy()
    })

    it('should display warning modal and delete the playlist when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/playlists/${playlist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const [_, store] = reduxRouterRender(withToastify(<PlaylistDrawer />), {
        global: {
          documentTitle: prevDocumentTitle,
          playlistDrawer: {
            open: true,
            playlistId: playlist.id
          },
          albumDrawer: undefined,
          artistDrawer: undefined,
          songDrawer: undefined
        }
      })

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete playlist/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete playlist/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.playlistDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.playlistDrawer.playlistId).toBeUndefined()
      expect((store.getState() as RootState).global.documentTitle).toBe(prevDocumentTitle)
      expect(screen.getByText(`${playlist.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate to song on song image click', async () => {
    const user = userEvent.setup()

    const song = songs[0]

    const [_, store] = render()

    await user.click(await screen.findByRole('img', { name: song.title }))
    expect((store.getState() as RootState).global.playlistDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/song/${song.id}`)
  })
})
