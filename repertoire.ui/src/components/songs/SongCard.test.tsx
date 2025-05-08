import { screen } from '@testing-library/react'
import SongCard from './SongCard'
import Song from '../../types/models/Song'
import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../test-utils'
import Artist from '../../types/models/Artist.ts'
import { userEvent } from '@testing-library/user-event'
import Difficulty from '../../types/enums/Difficulty.ts'
import { RootState } from '../../state/store.ts'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'

describe('Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '',
    title: 'Some song'
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render and display minimal info', () => {
    reduxRouterRender(<SongCard song={song} />)

    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByText(/unknown/i)).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(<SongCard song={localSong} />)

    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )

    const localSongWithAlbum: Song = {
      ...song,
      album: {
        ...emptyAlbum,
        imageUrl: 'something-album.png'
      }
    }

    rerender(<SongCard song={localSongWithAlbum} />)

    expect(screen.getByRole('img', { name: localSongWithAlbum.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  it('should render and display icons when the song is recorded, has guitar tuning, and songsterr and youtube links', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      artist: artist,
      isRecorded: true,
      guitarTuning: {
        id: '',
        name: 'Drop D'
      },
      difficulty: Difficulty.Impossible,
      songsterrLink: 'this is a link',
      youtubeLink: 'this is a link'
    }

    reduxRouterRender(<SongCard song={localSong} />)

    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByLabelText('recorded-icon')).toBeInTheDocument()
    expect(screen.getByLabelText('guitar-tuning-icon')).toBeInTheDocument()
    expect(screen.getByLabelText('difficulty-icon')).toBeInTheDocument()
    expect(screen.getByLabelText('songsterr-icon')).toBeInTheDocument()
    expect(screen.getByLabelText('youtube-icon')).toBeInTheDocument()

    await user.hover(screen.getByLabelText('recorded-icon'))
    expect(await screen.findByText(/is recorded/i)).toBeInTheDocument()

    await user.hover(screen.getByLabelText('guitar-tuning-icon'))
    expect(await screen.findByText(new RegExp(localSong.guitarTuning.name))).toBeInTheDocument()

    await user.hover(screen.getByLabelText('difficulty-icon'))
    expect(await screen.findByText(new RegExp(localSong.difficulty))).toBeInTheDocument()

    await user.hover(screen.getByLabelText('songsterr-icon'))
    expect(await screen.findByText(/songsterr/i)).toBeInTheDocument()

    await user.hover(screen.getByLabelText('youtube-icon'))
    expect(await screen.findByText(/youtube/i)).toBeInTheDocument()
  })

  it('should render and display solo icon when the song has exactly one Solo section', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      solosCount: 1
    }

    reduxRouterRender(<SongCard song={localSong} />)

    expect(screen.getByLabelText('solo-icon')).toBeInTheDocument()

    await user.hover(screen.getByLabelText('solo-icon'))
    expect(await screen.findByText(/a solo/i)).toBeInTheDocument()
  })

  it('should render and display solos icon when the song has more than one Solo section', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      solosCount: 2
    }

    reduxRouterRender(<SongCard song={localSong} />)

    expect(screen.getByLabelText('solos-icon')).toBeInTheDocument()

    await user.hover(screen.getByLabelText('solos-icon'))
    expect(await screen.findByText(/2 solos/i)).toBeInTheDocument()
  })

  it('should render and display riffs icon when the song has Riff sections', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      riffsCount: 2
    }

    reduxRouterRender(<SongCard song={localSong} />)

    expect(screen.getByLabelText('riffs-icon')).toBeInTheDocument()

    await user.hover(screen.getByLabelText('riffs-icon'))
    expect(await screen.findByText(/2 riffs/i)).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(<SongCard song={song} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`default-icon-${song.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /open drawer/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /view artist/i })).toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).toBeDisabled()

    rerender(<SongCard song={{ ...song, artist: emptyArtist, album: emptyAlbum }} />)
    expect(screen.getByRole('menuitem', { name: /view artist/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /view album/i })).not.toBeDisabled()
  })

  describe('on menu', () => {
    it('should open song drawer when clicking on open drawer', async () => {
      const user = userEvent.setup()

      const [_, store] = reduxRouterRender(withToastify(<SongCard song={song} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${song.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /open drawer/i }))

      expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
      expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
    })

    it('should navigate to artist when clicking on view artist', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, artist: { ...emptyArtist, id: '1' } }

      reduxRouterRender(<SongCard song={localSong} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localSong.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view artist/i }))
      expect(window.location.pathname).toBe(`/artist/${localSong.artist.id}`)
    })

    it('should navigate to album when clicking on view album', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, album: { ...emptyAlbum, id: '1' } }

      reduxRouterRender(<SongCard song={localSong} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${localSong.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /view album/i }))
      expect(window.location.pathname).toBe(`/album/${localSong.album.id}`)
    })

    it('should display warning modal when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<SongCard song={song} />))

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`default-icon-${song.title}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${song.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SongCard song={song} />)

    await user.click(screen.getByLabelText(`default-icon-${song.title}`))
    expect(window.location.pathname).toBe(`/song/${song.id}`)
  })

  it('should open artist drawer on artist name click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      artist: artist
    }

    const [_, store] = reduxRouterRender(<SongCard song={localSong} />)

    await user.click(screen.getByText(localSong.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })

  it('should open youtube modal on youtube icon click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      youtubeLink: 'https://youtube.com/watch?v=123456789'
    }

    server.use(
      http.get(localSong.youtubeLink.replace('watch?v=', 'embed/'), () => {
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    reduxRouterRender(<SongCard song={localSong} />)

    await user.click(screen.getByLabelText('youtube-icon'))

    expect(await screen.findByRole('dialog', { name: song.title })).toBeInTheDocument()
  })
})
