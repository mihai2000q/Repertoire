import Album from './Album.tsx'
import { reduxMemoryRouterRender } from '../test-utils.tsx'
import { screen } from '@testing-library/react'
import Song from '../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { default as AlbumType } from './../types/models/Album.ts'

describe('Album', () => {
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
      title: 'Song 1'
    },
    {
      ...emptySong,
      id: '2',
      title: 'Song 2'
    }
  ]

  const album: AlbumType = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
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
    })
  ]

  const getSongs = http.get('/songs', async () => {
    const response: WithTotalCountResponse<Song> = {
      models: songs,
      totalCount: songs.length
    }
    return HttpResponse.json(response)
  })

  const getAlbum = http.get(`/albums/${album.id}`, async () => {
    return HttpResponse.json(album)
  })

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render and display info from album when the album is not unknown', async () => {
    // Arrange
    server.use(getAlbum)

    // Act
    reduxMemoryRouterRender(<Album />, '/album/:id', [`/album/${album.id}`])

    // Assert
    expect(screen.getByTestId('album-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('songs-card')).toBeInTheDocument()
  })

  it('should render and display info from songs when the album is unknown', async () => {
    // Arrange
    server.use(getSongs)

    // Act
    reduxMemoryRouterRender(<Album />, '/album/:id', ['/album/unknown'])

    // Assert
    expect(screen.getByTestId('album-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('songs-card')).toBeInTheDocument()
  })
})
