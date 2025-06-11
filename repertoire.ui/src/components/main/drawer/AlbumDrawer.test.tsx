import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../../test-utils.tsx'
import AlbumDrawer from './AlbumDrawer.tsx'
import Album from '../../../types/models/Album.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { act, screen, waitFor } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import Artist from '../../../types/models/Artist.ts'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'
import dayjs from 'dayjs'
import { expect } from 'vitest'
import { openAlbumDrawer, setDocumentTitle } from '../../../state/slice/globalSlice.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../../types/models/Playlist.ts'

describe('Album Drawer', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      albumTrackNo: 1,
      imageUrl: 'something.png',
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      albumTrackNo: 2,
    },
    {
      ...emptySong,
      id: '3',
      title: 'Song 3',
      albumTrackNo: 3,
    },
  ]

  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1',
    songs: songs
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    imageUrl: 'something.png',
  }

  const getAlbum = (album: Album) =>
    http.get(`/albums/${album.id}`, async () => {
      return HttpResponse.json(album)
    })

  const handlers = [
    getAlbum(album),
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

  const render = (id: string | null = album.id) =>
    reduxRouterRender(<AlbumDrawer />, {
      global: {
        documentTitle: prevDocumentTitle,
        albumDrawer: {
          open: true,
          albumId: id
        },
        artistDrawer: undefined,
        songDrawer: undefined,
        playlistDrawer: undefined
      }
    })

  it('should render and display minimal info', async () => {
    const [_, store] = render()

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(`${album.songs.length} songs`)).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + album.title
    )

    album.songs.forEach((song) => {
      expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else {
        expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
      }
    })
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      imageUrl: 'something.png',
      releaseDate: '2024-10-16',
      artist: artist
    }

    server.use(getAlbum(localAlbum))

    const [_, store] = render()

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByRole('heading', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByText(dayjs(localAlbum.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toHaveAttribute(
      'src',
      localAlbum.artist.imageUrl
    )
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
    expect(screen.getByText(`${localAlbum.songs.length} songs`)).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(
      prevDocumentTitle + ' - ' + localAlbum.title
    )

    await user.hover(screen.getByText(dayjs(localAlbum.releaseDate).format('YYYY')))
    expect(
      await screen.findByText(new RegExp(dayjs(localAlbum.releaseDate).format('D MMMM YYYY')))
    ).toBeInTheDocument()

    localAlbum.songs.forEach((song) => {
      expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      if (song.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', song.imageUrl)
      } else if (album?.imageUrl) {
        expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
          'src',
          album.imageUrl
        )
      }
    })
  })

  it('should render the loader when there is no album', async () => {
    render(null)

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
  })

  it('should change document title on opening and closing, correctly', async () => {
    const newAlbum: Album = {
      ...album,
      id: 'another-id-another-album',
      title: 'new Title'
    }

    const [_, store] = render()

    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + album.title
      )
    })

    // change back the document title (as if the drawer closed)
    await act(() => store.dispatch(setDocumentTitle(prevDocumentTitle)))

    // make sure it doesn't use the old title when a new album is introduced
    server.use(getAlbum(newAlbum))
    await act(() => store.dispatch(openAlbumDrawer(newAlbum.id)))
    await waitFor(() => {
      expect((store.getState() as RootState).global.documentTitle).toBe(
        prevDocumentTitle + ' - ' + newAlbum.title
      )
    })
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      const [_, store] = render()

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/album/${album.id}`)
      expect((store.getState() as RootState).global.albumDrawer.open).toBeFalsy()
    })

    it('should display warning modal and delete the album when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const [_, store] = reduxRouterRender(withToastify(<AlbumDrawer />), {
        global: {
          documentTitle: prevDocumentTitle,
          albumDrawer: {
            open: true,
            albumId: album.id
          },
          artistDrawer: undefined,
          songDrawer: undefined,
          playlistDrawer: undefined
        }
      })

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.albumDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.albumDrawer.albumId).toBeUndefined()
      expect((store.getState() as RootState).global.documentTitle).toBe(prevDocumentTitle)
      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should navigate to artist on artist click', async () => {
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      artist: artist
    }

    server.use(getAlbum(localAlbum))

    const [_, store] = render()

    await user.click(await screen.findByText(localAlbum.artist.name))
    expect((store.getState() as RootState).global.albumDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/artist/${localAlbum.artist.id}`)
  })

  it('should navigate to song on song image click', async () => {
    const user = userEvent.setup()

    const song = songs[0]

    const [_, store] = render()

    await user.click(await screen.findByRole('img', { name: song.title }))
    expect((store.getState() as RootState).global.albumDrawer.open).toBeFalsy()
    expect(window.location.pathname).toBe(`/song/${song.id}`)
  })
})
