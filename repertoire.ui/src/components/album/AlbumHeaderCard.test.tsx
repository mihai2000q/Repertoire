import Album from 'src/types/models/Album.ts'
import {
  emptyAlbum,
  emptyArtist,
  emptySong,
  reduxRouterRender,
  withToastify
} from '../../test-utils.tsx'
import AlbumHeaderCard from './AlbumHeaderCard.tsx'
import userEvent from '@testing-library/user-event'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import Artist from 'src/types/models/Artist.ts'
import { RootState } from 'src/state/store.ts'
import dayjs from 'dayjs'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'

describe('Album Header Card', () => {
  const album: Album = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1',
    releaseDate: null
  }

  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    imageUrl: 'something.png'
  }

  const handlers = [
    http.get('/search', async () => {
      const response: WithTotalCountResponse<Artist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display minimal info when the album is not unknown', () => {
    reduxRouterRender(
      <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
    )

    expect(screen.getByLabelText(`default-icon-${album.title}`)).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText('0 songs')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()
  })

  it('should render and display maximal info when the album is not unknown', async () => {
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      imageUrl: 'something.png',
      releaseDate: '2024-10-21',
      artist: artist,
      songs: [
        {
          ...emptySong,
          id: '1',
          title: 'Song 1'
        },
        {
          ...emptySong,
          id: '2',
          title: 'Song 2'
        }
      ]
    }

    reduxRouterRender(
      <AlbumHeaderCard album={localAlbum} isUnknownAlbum={false} songsTotalCount={undefined} />
    )

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByRole('heading', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.artist.name })).toHaveAttribute(
      'src',
      localAlbum.artist.imageUrl
    )
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
    expect(screen.getByText(dayjs(localAlbum.releaseDate).year().toString())).toBeInTheDocument()
    expect(screen.getByText(`${localAlbum.songs.length} songs`)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'edit-header' })).toBeInTheDocument()

    await user.hover(screen.getByText(dayjs(localAlbum.releaseDate).year().toString()))
    expect(
      await screen.findByText(new RegExp(dayjs(localAlbum.releaseDate).format('D MMMM YYYY')))
    ).toBeInTheDocument()
  })

  it('should render and display info when the album is unknown', async () => {
    const songsTotalCount = 10

    reduxRouterRender(
      <AlbumHeaderCard album={undefined} isUnknownAlbum={true} songsTotalCount={songsTotalCount} />
    )

    expect(screen.getByLabelText('icon-unknown-album')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /unknown/i })).toBeInTheDocument()
    expect(screen.getByText(`${songsTotalCount} songs`)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'more-menu' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'edit-header' })).not.toBeInTheDocument()
  })

  it('should display image modal, when clicking the image', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumHeaderCard
        album={{ ...album, imageUrl: 'something.png' }}
        isUnknownAlbum={false}
        songsTotalCount={undefined}
      />
    )

    await user.click(screen.getByRole('img', { name: album.title }))
    expect(await screen.findByRole('dialog', { name: album.title + '-image' })).toBeInTheDocument()
  })

  it('should display menu on click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expect(screen.getByRole('menuitem', { name: /info/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should display info modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /info/i }))

      expect(await screen.findByRole('dialog', { name: /album info/i })).toBeInTheDocument()
    })

    it('should display edit header modal', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit album header/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete album', async () => {
      const user = userEvent.setup()

      let withSongs: string | null
      server.use(
        http.delete(`/albums/${album.id}`, (req) => {
          withSongs = new URL(req.request.url).searchParams.get('withSongs')
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(
        withToastify(
          <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(withSongs).toBe('false')

      expect(window.location.pathname).toBe('/albums')
      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
    })

    it('should display warning modal and delete album with songs', async () => {
      const user = userEvent.setup()

      let withSongs: string | null
      server.use(
        http.delete(`/albums/${album.id}`, (req) => {
          withSongs = new URL(req.request.url).searchParams.get('withSongs')
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRouterRender(
        withToastify(
          <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete album/i })).toBeInTheDocument()
      expect(screen.getByRole('checkbox', { name: /songs/i })).toBeInTheDocument()
      await user.click(screen.getByRole('checkbox', { name: /songs/i }))
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(withSongs).toBe('true')

      expect(window.location.pathname).toBe('/albums')
      expect(screen.getByText(`${album.title} deleted!`)).toBeInTheDocument()
    })
  })

  it('should display edit header modal from edit button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumHeaderCard album={album} isUnknownAlbum={false} songsTotalCount={undefined} />
    )

    await user.click(screen.getByRole('button', { name: 'edit-header' }))

    expect(await screen.findByRole('dialog', { name: /edit album header/i })).toBeInTheDocument()
  })

  it('should open artist drawer on artist click', async () => {
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      artist: artist
    }

    const [_, store] = reduxRouterRender(
      <AlbumHeaderCard album={localAlbum} isUnknownAlbum={false} songsTotalCount={undefined} />
    )

    await user.click(screen.getByText(localAlbum.artist.name))

    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })
})
