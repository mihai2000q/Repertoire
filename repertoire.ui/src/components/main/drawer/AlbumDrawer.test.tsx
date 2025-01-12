import {reduxRouterRender, withToastify} from '../../../test-utils.tsx'
import AlbumDrawer from './AlbumDrawer.tsx'
import Album from '../../../types/models/Album.ts'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { screen } from '@testing-library/react'
import Song from '../../../types/models/Song.ts'
import Artist from "../../../types/models/Artist.ts";
import {userEvent} from "@testing-library/user-event";
import {RootState} from "../../../state/store.ts";
import dayjs from "dayjs";

describe('Album Drawer', () => {
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

  const songs: Song[] = [
    {
      ...emptySong,
      id: '1',
      title: 'Song 1',
      albumTrackNo: 1
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2',
      albumTrackNo: 2
    }
  ]

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: songs
  }

  const artist: Artist = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const getAlbum = (album: Album) =>
    http.get(`/albums/${album.id}`, async () => {
      return HttpResponse.json(album)
    })

  const handlers = [getAlbum(album)]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  const render = (id: string | null = album.id) =>
    reduxRouterRender(<AlbumDrawer />, {
      global: {
        albumDrawer: {
          open: true,
          albumId: id
        },
        artistDrawer: undefined,
        songDrawer: undefined
      }
    })

  it('should render and display minimal info', async () => {
    render()

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(`${album.songs.length} songs`)).toBeInTheDocument()

    album.songs.forEach((song) => {
      expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
    })
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

    const localAlbum = {
      ...album,
      releaseDate: '2024-10-16',
      artist: artist
    }

    server.use(getAlbum(localAlbum))

    render()

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByText(dayjs(localAlbum.releaseDate).format('YYYY'))).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
    expect(screen.getByText(`${localAlbum.songs.length} songs`)).toBeInTheDocument()

    await user.hover(screen.getByText(dayjs(localAlbum.releaseDate).format('YYYY')))
    expect(await screen.findByText(new RegExp(dayjs(localAlbum.releaseDate).format('DD MMMM YYYY')))).toBeInTheDocument()

    localAlbum.songs.forEach((song) => {
      expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
      expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
      expect(screen.getByText(song.title)).toBeInTheDocument()
    })
  })

  it('should render the loader when there is no album', async () => {
    render(null)

    expect(screen.getByTestId('album-drawer-loader')).toBeInTheDocument()
  })

  it('should display menu when clicking on more', async () => {
    const user = userEvent.setup()

    render()

    await user.click(await screen.findByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to album when clicking on view details', async () => {
      const user = userEvent.setup()

      // Act
      render()

      // Assert
      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))
      expect(window.location.pathname).toBe(`/album/${album.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and delete the album when clicking delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/albums/${album.id}`, async () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      // Act
      const [_, store] = reduxRouterRender(withToastify(<AlbumDrawer />), {
        global: {
          albumDrawer: {
            open: true,
            albumId: album.id,
          },
          artistDrawer: undefined,
          songDrawer: undefined
        }
      })

      // Assert
      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

      expect((store.getState() as RootState).global.albumDrawer.open).toBeFalsy()
      expect((store.getState() as RootState).global.albumDrawer.albumId).toBeUndefined()
      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
    })
  })
})
