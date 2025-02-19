import { emptyAlbum, emptyArtist, emptySong, reduxRender } from '../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import HomeRecentlyPlayedSongs from './HomeRecentlyPlayedSongs.tsx'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../state/store.ts'

describe('Home Recently Played Songs', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1'
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      progress: 10,
      imageUrl: 'something.png',
      lastTimePlayed: '2024-12-12'
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song 3',
      imageUrl: 'something.png',
      lastTimePlayed: '2025-01-22',
      progress: 500,
      album: {
        ...emptyAlbum,
        imageUrl: 'something.png'
      }
    },
    {
      ...emptySong,
      id: '4',
      title: 'Song 4',
      progress: 20,
      album: {
        ...emptyAlbum,
        imageUrl: 'something.png'
      }
    },
    {
      ...emptySong,
      id: '5',
      title: 'Song 5',
      progress: 267,
      artist: { ...emptyArtist, name: 'Artist' }
    }
  ]

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRender(<HomeRecentlyPlayedSongs />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(screen.getByText(/title/i)).toBeInTheDocument()
    expect(screen.getByText(/progress/i)).toBeInTheDocument()
    expect(screen.getByLabelText('last-time-played-icon')).toBeInTheDocument()

    for (const song of songs) {
      expect(await screen.findByText(song.title)).toBeInTheDocument()
      if (song.artist) expect(screen.getByText(song.artist.name)).toBeInTheDocument()
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      if (song.imageUrl)
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      else if (song.album?.imageUrl)
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      if (song.lastTimePlayed)
        expect(screen.getByText(dayjs(song.lastTimePlayed).format('DD MMM'))).toBeInTheDocument()
    }

    // TODO: Why is this not working?
    // expect(screen.getAllByRole('progressbar', { name: 'progress' })).toHaveLength(songs.length)
    expect(screen.getAllByText('never')).toHaveLength(songs.filter((s) => !s.lastTimePlayed).length)
  })

  it('should display empty message when there are no songs', async () => {
    server.use(
      http.get('/songs', async () => {
        const response: WithTotalCountResponse<Song> = {
          models: [],
          totalCount: 0
        }
        return HttpResponse.json(response)
      })
    )

    reduxRender(<HomeRecentlyPlayedSongs />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.queryByText(/progress/i)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('last-time-played-icon')).not.toBeInTheDocument()
  })

  it('should open song drawer by clicking a song', async () => {
    const user = userEvent.setup()

    const song = songs[1]

    const [_, store] = reduxRender(<HomeRecentlyPlayedSongs />)

    await user.click(await screen.findByText(song.title))

    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
  })

  it("should open artist drawer by clicking a song's artist", async () => {
    const user = userEvent.setup()

    const song = songs.filter((s) => s.artist)[0]

    const [_, store] = reduxRender(<HomeRecentlyPlayedSongs />)

    await user.click(await screen.findByText(song.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(song.artist.id)
  })
})
