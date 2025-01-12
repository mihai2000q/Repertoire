import { reduxRouterRender, withToastify } from '../../../test-utils.tsx'
import SongDrawer from './SongDrawer.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import Artist from '../../../types/models/Artist.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import Album from '../../../types/models/Album.ts'
import Difficulty from '../../../utils/enums/Difficulty.ts'
import dayjs from 'dayjs'

describe('Song Drawer', () => {
  const song: Song = {
    id: '1',
    title: 'Song 1',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: ''
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    releaseDate: '2023-09-10',
    songs: []
  }

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const getSong = (song: Song) =>
    http.get(`/songs/${song.id}`, async () => {
      return HttpResponse.json(song)
    })

  const handlers = [getSong(song)]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const render = (id: string | null = song.id) =>
    reduxRouterRender(<SongDrawer />, {
      global: {
        songDrawer: {
          open: true,
          songId: id
        },
        artistDrawer: undefined,
        albumDrawer: undefined
      }
    })

  it('should render and display minimal info', async () => {
    render()

    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: song.title })).toBeInTheDocument()
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      description: 'This is a short description of the song',
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

    render()

    // Header
    expect(screen.getByTestId('song-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByText(dayjs(localSong.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()

    // hover release date
    await user.hover(screen.getByText(dayjs(localSong.releaseDate).format('YYYY')))
    expect(
      await screen.findByRole('tooltip', {
        name: new RegExp(dayjs(localSong.releaseDate).format('DD MMMM YYYY'))
      })
    ).toBeInTheDocument()

    // hover album
    await user.hover(screen.getByText(localSong.album.title))
    expect(await screen.findByRole('dialog', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getAllByText(localSong.album.title)).toHaveLength(2)
    expect(
      screen.getByText(new RegExp(dayjs(localSong.album.releaseDate).format('DD MMM YYYY')))
    ).toBeInTheDocument()

    // Details
    expect(screen.getByText(localSong.description)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    expect(screen.getByText(localSong.guitarTuning.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.bpm)).toBeInTheDocument()
    expect(screen.getByLabelText('recorded-icon')).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
    ).toBeInTheDocument()
    expect(screen.getByText(localSong.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(localSong.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(localSong.progress / 10)

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

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      render()

      // Assert
      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/song/${song.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and delete the song when clicking delete', async () => {
      // Arrange
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      const [_, store] = reduxRouterRender(withToastify(<SongDrawer />), {
        global: {
          songDrawer: {
            open: true,
            songId: song.id
          },
          artistDrawer: undefined,
          albumDrawer: undefined
        }
      })

      // Assert
      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('heading', { name: /delete song/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.songDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.songDrawer.songId).toBeUndefined()
      expect(screen.getByText(`${song.title} deleted!`)).toBeInTheDocument()
    })
  })
})
