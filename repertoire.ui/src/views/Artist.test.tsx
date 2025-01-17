import Artist from './Artist.tsx'
import { reduxMemoryRouterRender } from '../test-utils.tsx'
import { screen } from '@testing-library/react'
import Song from '../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { default as ArtistType } from './../types/models/Artist.ts'

describe('Artist', () => {
  const artist: ArtistType = {
    id: '1',
    name: 'Artist 1',
    createdAt: '',
    updatedAt: '',
    albums: [],
    songs: []
  }

  const albumResponse: WithTotalCountResponse<Song> = {
    models: [],
    totalCount: 0
  }
  const songsResponse: WithTotalCountResponse<Song> = {
    models: [],
    totalCount: 0
  }

  const getArtist = http.get(`/artists/${artist.id}`, async () => {
    return HttpResponse.json(artist)
  })

  const server = setupServer()

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display artist info when the artist is not unknown', async () => {
    let albumsParams: URLSearchParams
    let songsParams: URLSearchParams
    server.use(
      getArtist,
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

    reduxMemoryRouterRender(<Artist />, '/artist/:id', [`/artist/${artist.id}`])

    expect(screen.getByTestId('artist-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('albums-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('songs-card')).toBeInTheDocument()

    expect(albumsParams.getAll('searchBy')).toStrictEqual([`artist_id = '${artist.id}'`])
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

    reduxMemoryRouterRender(<Artist />, '/artist/:id', ['/artist/unknown'])

    expect(screen.getByLabelText('header-panel-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('albums-card')).toBeInTheDocument()
    expect(await screen.findByLabelText('songs-card')).toBeInTheDocument()

    expect(albumsParams.getAll('searchBy')).toStrictEqual([`artist_id IS NULL`])
    expect(songsParams.getAll('searchBy')).toStrictEqual([`songs.artist_id IS NULL`])
  })
})
