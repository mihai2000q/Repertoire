import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../../test-utils.tsx'
import SongDrawer from './SongDrawer.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { act, screen, waitFor, within } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import Artist from '../../../types/models/Artist.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import Album from '../../../types/models/Album.ts'
import Difficulty from '../../../types/enums/Difficulty.ts'
import dayjs from 'dayjs'
import { expect } from 'vitest'
import {
  closeSongDrawer,
  openSongDrawer,
  setDocumentTitle
} from '../../../state/slice/globalSlice.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../../types/models/Playlist.ts'

describe('Song Drawer', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1'
  }

  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1',
    releaseDate: '2023-09-10',
    imageUrl: 'something2.png'
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    imageUrl: 'something3.png'
  }

  const getSong = (song: Song) =>
    http.get(`/songs/${song.id}`, async () => {
      return HttpResponse.json(song)
    })

  const handlers = [
    getSong(song),
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

  const render = (id: string | null = song.id) =>
    reduxRouterRender(<SongDrawer />, {
      global: {
        documentTitle: prevDocumentTitle,
        songDrawer: {
          open: true,
          songId: id
        },
        artistDrawer: undefined,
        albumDrawer: undefined,
        playlistDrawer: undefined
      }
    })

  it('should render and display minimal info', async () => {
    const [_, store] = render()

    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: song.title })).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + song.title
    )
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      description: 'This is a short description of the song',
      imageUrl: 'something.png',
      releaseDate: '2024-10-16',
      difficulty: Difficulty.Medium,
      isRecorded: true,
      bpm: 120,
      guitarTuning: {
        id: '1',
        name: 'Drop D'
      },
      artist: artist,
      album: album,
      songsterrLink: 'some link',
      youtubeLink: 'some other link',
      lastTimePlayed: '2023-10-16',
      rehearsals: 15,
      confidence: 50,
      progress: 210
    }

    server.use(getSong(localSong))

    const [_, store] = render()

    // Header
    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByRole('heading', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByText(dayjs(localSong.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toHaveAttribute(
      'src',
      localSong.artist.imageUrl
    )
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + localSong.title
    )

    // hover release date
    await user.hover(screen.getByText(dayjs(localSong.releaseDate).format('YYYY')))
    expect(
      await screen.findByRole('tooltip', {
        name: new RegExp(dayjs(localSong.releaseDate).format('D MMMM YYYY'))
      })
    ).toBeInTheDocument()

    // hover album
    await user.hover(screen.getByText(localSong.album.title))
    expect(await screen.findByRole('dialog', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.album.title })).toHaveAttribute(
      'src',
      localSong.album.imageUrl
    )
    expect(screen.getAllByText(localSong.album.title)).toHaveLength(2)
    expect(
      screen.getByText(new RegExp(dayjs(localSong.album.releaseDate).format('D MMM YYYY')))
    ).toBeInTheDocument()

    // Details
    expect(screen.getByText(localSong.description)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    expect(screen.getByText(localSong.guitarTuning.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.bpm)).toBeInTheDocument()
    expect(screen.getByLabelText('recorded-icon')).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(localSong.lastTimePlayed).format('D MMMM YYYY'))
    ).toBeInTheDocument()
    expect(screen.getByText(localSong.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(
      localSong.confidence
    )
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(
      localSong.progress / 10
    )

    // hover difficulty, confidence and progress bars
    await user.hover(screen.getByRole('progressbar', { name: 'difficulty' }))
    expect(
      await screen.findByRole('tooltip', { name: new RegExp(localSong.difficulty) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(localSong.confidence.toString()) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(localSong.progress.toString()) })
    ).toBeInTheDocument()

    // links
    expect(screen.getByRole('button', { name: 'songsterr' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'youtube' })).toBeInTheDocument()
  })

  it('should render the loader when there is no song', async () => {
    render(null)

    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", async () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    server.use(getSong(localSong))

    const [{ rerender }, store] = render(localSong.id)

    expect(await screen.findByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )

    const localSongWithAlbum: Song = {
      ...song,
      id: '3123123123',
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    }

    server.use(getSong(localSongWithAlbum))
    await act(() => store.dispatch(closeSongDrawer()))
    await act(() => store.dispatch(openSongDrawer(localSongWithAlbum.id)))
    rerender(<SongDrawer />)

    await waitFor(() => {
      expect(screen.getByRole('img', { name: localSongWithAlbum.title })).toHaveAttribute(
        'src',
        localSongWithAlbum.album.imageUrl
      )
    })
  })

  it('should change document title on opening and closing, correctly', async () => {
    const newSong: Song = {
      ...song,
      id: 'another-id-another-song',
      title: 'new Title'
    }

    const [_, store] = render()

    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + song.title
      )
    })

    // change back the document title (as if the drawer closed)
    await act(() => store.dispatch(setDocumentTitle(prevDocumentTitle)))

    // make sure it doesn't use the old title when a new album is introduced
    server.use(getSong(newSong))
    await act(() => store.dispatch(openSongDrawer(newSong.id)))
    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + newSong.title
      )
    })
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      const [_, store] = render()

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect((store.getState() as RootState).global.songDrawer.open).toBeFalsy()
      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should display warning modal and delete the song when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const [_, store] = reduxRouterRender(withToastify(<SongDrawer />), {
        global: {
          documentTitle: prevDocumentTitle,
          songDrawer: {
            open: true,
            songId: song.id
          },
          artistDrawer: undefined,
          albumDrawer: undefined,
          playlistDrawer: undefined
        }
      })

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete song/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.songDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.songDrawer.songId).toBeUndefined()
      expect((store.getState() as RootState).global.documentTitle).toBe(prevDocumentTitle)
      expect(screen.getByText(`${song.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate to artist on artist click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      artist: artist
    }

    server.use(getSong(localSong))

    const [_, store] = render()

    await user.click(await screen.findByText(localSong.artist.name))
    expect((store.getState() as RootState).global.songDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/artist/${localSong.artist.id}`)
  })

  it('should navigate to album on album click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      album: album
    }

    server.use(getSong(localSong))

    const [_, store] = render()

    await user.click(await screen.findByText(localSong.album.title))
    expect((store.getState() as RootState).global.songDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/album/${localSong.album.id}`)
  })

  it('should open youtube modal on youtube icon click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      youtubeLink: 'https://www.youtube.com/watch?v=E0ozmU9cJDg'
    }

    server.use(
      getSong(localSong),
      http.get(
        localSong.youtubeLink.replace('youtube', 'youtube-nocookie').replace('watch?v=', 'embed/'),
        () => {
          return HttpResponse.json({ message: 'it worked' })
        }
      )
    )

    render()

    await user.click(await screen.findByRole('button', { name: 'youtube' }))

    expect(await screen.findByRole('dialog', { name: song.title })).toBeInTheDocument()
  })

  it('should be able to open songsterr in browser on songsterr click', async () => {
    const localSong = {
      ...song,
      songsterrLink: 'some-link'
    }

    server.use(getSong(localSong))

    render()

    expect(
      within(await screen.findByRole('link', { name: /songsterr/i })).getByRole('button', {
        name: /songsterr/i
      })
    ).toBeInTheDocument()

    expect(screen.getByRole('link', { name: /songsterr/i })).toBeExternalLink(
      localSong.songsterrLink
    )
  })
})
