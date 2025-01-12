import { reduxRouterRender, withToastify } from '../../../test-utils.tsx'
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

describe('Song Header Card', () => {
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
    updatedAt: '',
    releaseDate: null
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    releaseDate: '2022-10-15',
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

  const handlers = [
    http.get('/songs/guitar-tunings', async () => {
      return HttpResponse.json([])
    }),
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json([])
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
      releaseDate: '2024-10-21T10:30:00',
      artist: artist,
      album: album
    }

    reduxRouterRender(<SongHeaderCard song={localSong} />)

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.artist.name })).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localSong.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.hover(screen.getByText(dayjs(localSong.releaseDate).format('YYYY')))
    expect(
      await screen.findByText(new RegExp(dayjs(localSong.releaseDate).format('DD MMMM YYYY')))
    ).toBeInTheDocument()

    await user.hover(screen.getByText(localSong.album.title))
    expect(await screen.findByRole('img', { name: localSong.album.title })).toBeInTheDocument()
    expect(screen.getAllByText(localSong.album.title)).toHaveLength(2)
    expect(
      screen.getByText(dayjs(localSong.album.releaseDate).format('DD MMM YYYY'))
    ).toBeInTheDocument()
  })

  it('should display menu on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(<SongHeaderCard song={song} />)

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<SongHeaderCard song={song} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(screen.getByRole('heading', { name: /song info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      // Arrange
      const user = userEvent.setup()

      // Act
      reduxRouterRender(<SongHeaderCard song={song} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(screen.getByRole('heading', { name: /edit song header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete song', async () => {
      // Arrange
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      reduxRouterRender(withToastify(<SongHeaderCard song={song} />))

      // Assert
      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

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

    expect(screen.getByRole('heading', { name: /edit song header/i })).toBeInTheDocument()
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
