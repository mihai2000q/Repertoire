import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../../test-utils.tsx'
import SongHeaderCard from './SongHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import Song from 'src/types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import Artist from 'src/types/models/Artist.ts'
import { RootState } from 'src/state/store.ts'
import Album from '../../../types/models/Album.ts'
import dayjs from 'dayjs'
import WithTotalCountResponse from "../../../types/responses/WithTotalCountResponse.ts";

describe('Song Header Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1'
  }

  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1',
    releaseDate: '2022-10-15'
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const handlers = [
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json([])
    }),
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json([])
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get('/artists', async () => {
      const response: WithTotalCountResponse<Artist> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display minimal info', async () => {
    reduxRouterRender(<SongHeaderCard song={song} />)

    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: song.title })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      imageUrl: 'something.png',
      releaseDate: '2024-10-21T10:30:00',
      artist: artist,
      album: album
    }

    reduxRouterRender(<SongHeaderCard song={localSong} />)

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByRole('heading', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toHaveAttribute(
      'src',
      localSong.artist.imageUrl
    )
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localSong.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.hover(screen.getByText(dayjs(localSong.releaseDate).format('YYYY')))
    expect(
      await screen.findByText(new RegExp(dayjs(localSong.releaseDate).format('D MMMM YYYY')))
    ).toBeInTheDocument()

    await user.hover(screen.getByText(localSong.album.title))
    expect(await screen.findByRole('img', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.album.title })).toHaveAttribute(
      'src',
      localSong.album.imageUrl
    )
    expect(screen.getAllByText(localSong.album.title)).toHaveLength(2)
    expect(
      screen.getByText(dayjs(localSong.album.releaseDate).format('D MMM YYYY'))
    ).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(<SongHeaderCard song={localSong} />)

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', localSong.imageUrl)

    const localSongWithAlbum: Song = {
      ...song,
      album: {
        id: '',
        title: '',
        songs: [],
        createdAt: '',
        updatedAt: '',
        imageUrl: 'something-album.png'
      }
    }

    rerender(<SongHeaderCard song={localSongWithAlbum} />)

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  it('should display image modal, when clicking the image', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SongHeaderCard song={{ ...song, imageUrl: 'something.png' }} />)

    await user.click(screen.getByRole('img', { name: song.title }))
    expect(await screen.findByRole('dialog', { name: song.title + '-image' })).toBeInTheDocument()
  })

  it('should display menu on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SongHeaderCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<SongHeaderCard song={song} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(await screen.findByRole('dialog', { name: /song info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(<SongHeaderCard song={song} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit song header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete song', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(withToastify(<SongHeaderCard song={song} />))

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete song/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(window.location.pathname).toBe('/songs')
      expect(screen.getByText(`${song.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should display edit header modal from edit button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SongHeaderCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(await screen.findByRole('dialog', { name: /edit song header/i })).toBeInTheDocument()
  })

  it('should open artist drawer on artist click', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      artist: artist
    }

    const [_, store] = reduxRouterRender(<SongHeaderCard song={localSong} />)

    await user.click(screen.getByText(localSong.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })

  it('should open album drawer on album click', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      album: album
    }

    const [_, store] = reduxRouterRender(<SongHeaderCard song={localSong} />)

    await user.click(screen.getByText(localSong.album.title))

    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(localSong.album.id)
  })
})
