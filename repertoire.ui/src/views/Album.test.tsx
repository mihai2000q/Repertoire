import Album from './Album.tsx'
import { emptyAlbum, emptySong, reduxMemoryRouterRender } from '../test-utils.tsx'
import { screen } from '@testing-library/react'
import Song from '../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { default as AlbumType } from './../types/models/Album.ts'
import { RootState } from '../state/store.ts'
import { expect } from 'vitest'
import albumSongsOrders from '../data/album/albumSongsOrders.ts'
import Artist from '../types/models/Artist.ts'

describe('Album', () => {
  const songs: Song[] = [
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

  const album: AlbumType = {
    ...emptyAlbum,
    id: '1',
    title: 'Album 1',
    releaseDate: null,
    songs: [
      ...songs,
      {
        ...emptySong,
        id: '3',
        title: 'Song 3'
      }
    ]
  }

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: songs,
        totalCount: songs.length
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

  it('should render and display info from album when the album is not unknown', async () => {
    let songsOrderBy: string[]
    server.use(
      http.get(`/albums/${album.id}`, async (req) => {
        if (!songsOrderBy) {
          songsOrderBy = new URL(req.request.url).searchParams.getAll('songsOrderBy')
        }
        return HttpResponse.json(album)
      })
    )

    const [_, store] = reduxMemoryRouterRender(<Album />, '/album/:id', [`/album/${album.id}`])

    expect(screen.getByTestId('album-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('songs-card')).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(album.title)

    expect(songsOrderBy).toHaveLength(1)
    expect(songsOrderBy[0]).toBe(albumSongsOrders[0].value)
  })

  it('should render and display info from songs when the album is unknown', async () => {
    let songsSearchBy: string[]
    let songsOrderBy: string[]
    server.use(
      http.get('/songs', async (req) => {
        if (!songsSearchBy || !songsOrderBy) {
          songsSearchBy = new URL(req.request.url).searchParams.getAll('searchBy')
          songsOrderBy = new URL(req.request.url).searchParams.getAll('orderBy')
        }
        const response: WithTotalCountResponse<Song> = {
          models: songs,
          totalCount: songs.length
        }
        return HttpResponse.json(response)
      })
    )

    const [_, store] = reduxMemoryRouterRender(<Album />, '/album/:id', ['/album/unknown'])

    expect(screen.getByTestId('album-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('songs-card')).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toMatch(/unknown/i)

    expect(songsSearchBy).toHaveLength(1)
    expect(songsSearchBy[0]).toMatch('album_id IS NULL')

    expect(songsOrderBy).toHaveLength(1)
    expect(songsOrderBy[0]).toBe(albumSongsOrders[1].value)
  })
})
