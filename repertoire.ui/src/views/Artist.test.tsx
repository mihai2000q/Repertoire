import Artist from './Artist.tsx'
import { emptyArtist, reduxMemoryRouterRender } from '../test-utils.tsx'
import { screen } from '@testing-library/react'
import Song from '../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { default as ArtistType } from './../types/models/Artist.ts'
import { expect } from 'vitest'
import { RootState } from '../state/store.ts'
import Album from '../types/models/Album.ts'
import { SearchBase } from '../types/models/Search.ts'
import createOrder from '../utils/createOrder.ts'
import artistAlbumsOrders from '../data/artist/artistAlbumsOrders.ts'
import artistSongsOrders from '../data/artist/artistSongsOrders.ts'

describe('Artist', () => {
  const artist: ArtistType = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1'
  }

  const albumResponse: WithTotalCountResponse<Song> = {
    models: [],
    totalCount: 0
  }
  const songsResponse: WithTotalCountResponse<Song> = {
    models: [],
    totalCount: 0
  }

  const getArtist = (a: ArtistType = artist) =>
    http.get(`/artists/${a.id}`, async () => {
      return HttpResponse.json(a)
    })

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get('/search', async () => {
      const response: WithTotalCountResponse<SearchBase> = {
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

  it('should render and display artist info when the artist is not unknown', async () => {
    let albumsParams: URLSearchParams
    let songsParams: URLSearchParams
    server.use(
      getArtist(),
      http.get('/albums', async (req) => {
        if (!albumsParams) {
          albumsParams = new URL(req.request.url).searchParams
        }
        return HttpResponse.json(albumResponse)
      }),
      http.get('/songs', async (req) => {
        if (!songsParams) {
          songsParams = new URL(req.request.url).searchParams
        }
        return HttpResponse.json(songsResponse)
      })
    )

    const [_, store] = reduxMemoryRouterRender(<Artist />, '/artist/:id', [`/artist/${artist.id}`])

    expect(screen.getByTestId('artist-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('albums-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('songs-card')).toBeInTheDocument()
    expect(screen.queryByLabelText('band-members-card')).not.toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(artist.name)

    expect(albumsParams.getAll('orderBy')).toStrictEqual([
      createOrder(artistAlbumsOrders[0]),
      createOrder(artistAlbumsOrders[0].thenBy[0])
    ])
    expect(albumsParams.getAll('searchBy')).toStrictEqual([`artist_id = '${artist.id}'`])
    expect(albumsParams.getAll('orderBy')).toStrictEqual([
      createOrder(artistSongsOrders[0]),
      createOrder(artistSongsOrders[0].thenBy[0])
    ])
    expect(songsParams.getAll('searchBy')).toStrictEqual([`songs.artist_id = '${artist.id}'`])
  })

  it('should render and display only songs and albums when the artist is unknown', async () => {
    let albumsParams: URLSearchParams
    let songsParams: URLSearchParams
    server.use(
      http.get('/albums', async (req) => {
        if (!albumsParams) {
          albumsParams = new URL(req.request.url).searchParams
        }
        return HttpResponse.json(albumResponse)
      }),
      http.get('/songs', async (req) => {
        if (!songsParams) {
          songsParams = new URL(req.request.url).searchParams
        }
        return HttpResponse.json(songsResponse)
      })
    )

    const [_, store] = reduxMemoryRouterRender(<Artist />, '/artist/:id', ['/artist/unknown'])

    expect((store.getState() as RootState).global.documentTitle).toMatch(/unknown/i)
    expect(screen.getByLabelText('header-panel-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('albums-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('songs-card')).toBeInTheDocument()
    expect(screen.queryByLabelText('band-members-card')).not.toBeInTheDocument()

    expect(albumsParams.getAll('orderBy')).toStrictEqual([
      createOrder(artistAlbumsOrders[0]),
      createOrder(artistAlbumsOrders[0].thenBy[0])
    ])
    expect(albumsParams.getAll('searchBy')).toStrictEqual([`artist_id IS NULL`])
    expect(albumsParams.getAll('orderBy')).toStrictEqual([
      createOrder(artistSongsOrders[0]),
      createOrder(artistSongsOrders[0].thenBy[0])
    ])
    expect(songsParams.getAll('searchBy')).toStrictEqual([`songs.artist_id IS NULL`])
  })

  it('should render and display band members when the artist is a band', async () => {
    server.use(getArtist({ ...artist, isBand: true }))

    reduxMemoryRouterRender(<Artist />, '/artist/:id', [`/artist/${artist.id}`])

    expect(await screen.findByLabelText('band-members-card')).toBeInTheDocument()
  })
})
