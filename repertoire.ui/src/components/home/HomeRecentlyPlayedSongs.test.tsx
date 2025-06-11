import { emptyAlbum, emptyArtist, emptySong, reduxRouterRender } from '../../test-utils.tsx'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Song from '../../types/models/Song.ts'
import { screen, waitFor } from '@testing-library/react'
import HomeRecentlyPlayedSongs from './HomeRecentlyPlayedSongs.tsx'
import dayjs from 'dayjs'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../state/store.ts'
import { expect } from 'vitest'

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
        id: '1',
        imageUrl: 'something.png'
      },
      artist: {
        ...emptyArtist,
        id: '1',
        name: 'another-artist'
      },
      songsterrLink: 'some-link'
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

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render', async () => {
    reduxRouterRender(<HomeRecentlyPlayedSongs />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(screen.getByText(/title/i)).toBeInTheDocument()
    expect(screen.getByText(/progress/i)).toBeInTheDocument()
    expect(screen.getByLabelText('last-time-played-icon')).toBeInTheDocument()

    for (const song of songs) {
      expect(await screen.findByText(song.title)).toBeInTheDocument()
      if (song.artist) expect(screen.getByText(song.artist.name)).toBeInTheDocument()
      if (song.imageUrl)
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      else if (song.album?.imageUrl)
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      else expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
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

    reduxRouterRender(<HomeRecentlyPlayedSongs />)

    expect(screen.getByText(/recently played/i)).toBeInTheDocument()

    expect(await screen.findByText(/no songs/i)).toBeInTheDocument()

    expect(screen.queryByText(/title/i)).not.toBeInTheDocument()
    expect(screen.queryByText(/progress/i)).not.toBeInTheDocument()
    expect(screen.queryByLabelText('last-time-played-icon')).not.toBeInTheDocument()
  })

  it('should open song drawer by clicking a song', async () => {
    const user = userEvent.setup()

    const song = songs[1]

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedSongs />)

    await user.click(await screen.findByText(song.title))

    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
  })

  it("should open artist drawer by clicking a song's artist", async () => {
    const user = userEvent.setup()

    const song = songs.filter((s) => s.artist)[0]

    const [_, store] = reduxRouterRender(<HomeRecentlyPlayedSongs />)

    await user.click(await screen.findByText(song.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(song.artist.id)
  })

  it('should display menu on song right click', async () => {
    const user = userEvent.setup()

    const song1 = songs[0]
    const song2 = songs[2]

    const [{ rerender }] = reduxRouterRender(<HomeRecentlyPlayedSongs />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByLabelText(`default-icon-${song1.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()

    rerender(<HomeRecentlyPlayedSongs />)
    await user.pointer({
      keys: '[MouseRight>]',
      target: await screen.findByRole('img', { name: song2.title })
    })
    // the other menu is still opened,
    // waiting for it to close automatically by opening another one
    await waitFor(() =>
      expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    )
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      const song = songs[0]

      reduxRouterRender(<HomeRecentlyPlayedSongs />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByLabelText(`default-icon-${song.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const song = songs[2]

      reduxRouterRender(<HomeRecentlyPlayedSongs />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: song.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${song.artist.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      const song = songs[2]

      reduxRouterRender(<HomeRecentlyPlayedSongs />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: await screen.findByRole('img', { name: song.title })
      })
      await user.click(screen.getByRole('menuitem', { name: /view album/i }))
      expect(window.location.pathname).toBe(`/album/${song.album.id}`)
    })
  })
})
