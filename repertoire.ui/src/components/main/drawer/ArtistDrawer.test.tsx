import { emptyAlbum, emptyArtist, emptySong, reduxRouterRender } from '../../../test-utils.tsx'
import ArtistDrawer from './ArtistDrawer.tsx'
import Artist from '../../../types/models/Artist.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { act, screen, waitFor } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import Album from '../../../types/models/Album.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import dayjs from 'dayjs'
import { expect } from 'vitest'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import Playlist from '../../../types/models/Playlist.ts'

describe('Artist Drawer', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      imageUrl: 'something.png'
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2'
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song 3',
      imageUrl: 'something.png',
      album: {
        ...emptyAlbum,
        id: '1',
        title: 'Song Album 1',
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      id: '4',
      title: 'Song 4',
      album: {
        ...emptyAlbum,
        id: '2',
        title: 'Song Album 2',
        imageUrl: 'something-album.png'
      }
    },
    {
      ...emptySong,
      id: '5',
      title: 'Song 5',
      imageUrl: 'something.png',
      album: {
        ...emptyAlbum,
        id: '3',
        title: 'Song Album 3'
      }
    },
    {
      ...emptySong,
      id: '6',
      title: 'Song 6',
      album: {
        ...emptyAlbum,
        id: '4',
        title: 'Song Album 4'
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
      title: 'Album 2',
      imageUrl: 'something.png'
    }
  ]

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    imageUrl: 'something.png',
    isBand: true,
    bandMembers: [
      {
        id: '1',
        name: 'Member 1',
        roles: [],
        imageUrl: 'something.png'
      },
      {
        id: '2',
        name: 'Member 2',
        roles: []
      }
    ]
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
    }),
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

  const render = (id: string | null = artist.id) =>
    reduxRouterRender(<ArtistDrawer />, {
      global: {
        documentTitle: prevDocumentTitle,
        artistDrawer: {
          open: true,
          artistId: id
        },
        albumDrawer: undefined,
        songDrawer: undefined,
        playlistDrawer: undefined
      }
    })

  it('should render', async () => {
    const [_, store] = render()

    expect(screen.getByTestId('artist-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: artist.name })).toHaveAttribute('src', artist.imageUrl)
    expect(screen.getByRole('heading', { name: artist.name })).toBeInTheDocument()
    expect(
      screen.getByText(
        `${artist.bandMembers.length} members • ${albums.length} albums • ${songs.length} songs`
      )
    ).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + artist.name
    )

    artist.bandMembers.forEach((bandMember) => {
      expect(screen.getByText(bandMember.name)).toBeInTheDocument()
      if (bandMember.imageUrl) {
        expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
      } else {
        expect(screen.getByLabelText(`icon-${bandMember.name}`)).toBeInTheDocument()
      }
    })

    albums.forEach((album) => {
      expect(screen.getByText(album.title)).toBeInTheDocument()
      if (album.imageUrl) {
        expect(screen.getByRole('img', { name: album.title })).toHaveAttribute(
          'src',
          album.imageUrl
        )
      } else {
        expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
      }
      if (album.releaseDate) {
        expect(screen.getByText(dayjs(album.releaseDate).format('D MMM YYYY'))).toBeInTheDocument()
      }
    })
    songs.forEach((song) => {
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else if (song.album?.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          song.album.imageUrl
        )
      } else {
        expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
      }
      if (song.album) {
        expect(screen.getByText(song.album.title)).toBeInTheDocument()
      }
    })
  })

  it('should render the loader when there is no artist', async () => {
    render(null)

    expect(screen.getByTestId('artist-drawer-loader')).toBeInTheDocument()
  })

  it('should change document title on opening and closing, correctly', async () => {
    const newArtist: Artist = {
      ...artist,
      id: 'another-id-another-artist',
      name: 'new Name'
    }

    const [_, store] = render()

    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + artist.name
      )
    })

    // click outside to close drawer
    await userEvent.click(document.querySelector('.mantine-Drawer-overlay'))

    // make sure it doesn't use the old name when a new artist is introduced
    server.use(getArtist(newArtist))
    await act(() => store.dispatch(openArtistDrawer(newArtist.id)))
    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + newArtist.name
      )
    })
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to artist when clicking on view details', async () => {
      const user = userEvent.setup()

      const [_, store] = render()

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/artist/${artist.id}`)
      expect((store.getState() as RootState).global.artistDrawer.open).toBeFalsy()
    })

    it('should display warning modal and delete the artist when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/artists/${artist.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const [_, store] = reduxRouterRender(<ArtistDrawer />, {
        global: {
          documentTitle: prevDocumentTitle,
          artistDrawer: {
            open: true,
            artistId: artist.id
          },
          albumDrawer: undefined,
          songDrawer: undefined,
          playlistDrawer: undefined
        }
      })

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete artist/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.artistDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.artistDrawer.artistId).toBeUndefined()
      expect((store.getState() as RootState).global.documentTitle).toBe(prevDocumentTitle)
    })
  })

  it('should navigate to album on album image click', async () => {
    const user = userEvent.setup()

    const album = albums[1]

    const [_, store] = render()

    await user.click(await screen.findByRole('img', { name: album.title }))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/album/${album.id}`)
  })

  it('should navigate to song on song image click', async () => {
    const user = userEvent.setup()

    const song = songs[0]

    const [_, store] = render()

    await user.click(await screen.findByRole('img', { name: song.title }))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/song/${song.id}`)
  })
})
