import { reduxRouterRender, withToastify } from '../../../test-utils.tsx'
import ArtistDrawer from './ArtistDrawer.tsx'
import Artist from '../../../types/models/Artist.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import Album from '../../../types/models/Album.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import dayjs from 'dayjs'

describe('Artist Drawer', () => {
  const emptySong: Song = {
    id: '',
    title: '',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: ''
  }

  const emptyAlbum: Album = {
    id: '',
    title: '',
    songs: [],
    createdAt: '',
    updatedAt: ''
  }

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
      album: {
        ...emptyAlbum,
        id: '12',
        title: 'Song Album'
      }
    }
  ]

  const albums: Album[] = [
    {
      ...emptyAlbum,
      id: '1',
      title: 'Album 1',
      releaseDate: '2017-05-01'
    },
    {
      ...emptyAlbum,
      id: '2',
      title: 'Album 2'
    }
  ]

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const getArtist = (artist: Artist) =>
    http.get(`/artists/${artist.id}`, async () => {
      return HttpResponse.json(artist)
    })

  const handlers = [
    getArtist(artist),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: albums,
        totalCount: albums.length
      }
      return HttpResponse.json(response)
    }),
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

  const render = (id: string | null = artist.id) =>
    reduxRouterRender(<ArtistDrawer />, {
      global: {
        artistDrawer: {
          open: true,
          artistId: id
        },
        albumDrawer: undefined,
        songDrawer: undefined
      }
    })

  it('should render', async () => {
    render()

    expect(screen.getByTestId('artist-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByText(`${albums.length} albums • ${songs.length} songs`)).toBeInTheDocument()

    albums.forEach((album) => {
      expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
      expect(screen.getByText(album.title)).toBeInTheDocument()
      if (album.releaseDate) {
        expect(screen.getByText(dayjs(album.releaseDate).format('D MMM YYYY'))).toBeInTheDocument()
      }
    })
    songs.forEach((song) => {
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.album) {
        expect(screen.getByText(song.album.title)).toBeInTheDocument()
      }
    })
  })

  it('should render the loader when there is no artist', async () => {
    render(null)

    expect(screen.getByTestId('artist-drawer-loader')).toBeInTheDocument()
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to artist when clicking on view details', async () => {
      const user = userEvent.setup()

      render()

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/artist/${artist.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and delete the artist when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const [_, store] = reduxRouterRender(withToastify(<ArtistDrawer />), {
        global: {
          artistDrawer: {
            open: true,
            artistId: artist.id
          },
          albumDrawer: undefined,
          songDrawer: undefined
        }
      })

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('dialog', { name: /delete artist/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete artist/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.artistDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.artistDrawer.artistId).toBeUndefined()
      expect(screen.getByText(`${artist.name} deleted!`)).toBeInTheDocument()
    })
  })
})
