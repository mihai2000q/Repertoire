import { reduxRender } from '../test-utils.tsx'
import Home from './Home.tsx'
import { screen } from '@testing-library/react'
import { expect } from 'vitest'
import { RootState } from '../state/store.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import Song from '../types/models/Song.ts'
import Album from '../types/models/Album.ts'
import Artist from '../types/models/Artist.ts'
import { setupServer } from 'msw/node'
import Playlist from '../types/models/Playlist.ts'

describe('Home', () => {
  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/albums', async () => {
      const response: WithTotalCountResponse<Album> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get('/artists', async () => {
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

  it('should render', () => {
    const [_, store] = reduxRender(<Home />)

    expect((store.getState() as RootState).global.documentTitle).toMatch(/home/i)
    expect(screen.getByLabelText('top')).toBeInTheDocument()
    expect(screen.getByLabelText('genres')).toBeInTheDocument()
    expect(screen.getByLabelText('recently-played-songs')).toBeInTheDocument()
    expect(screen.getByLabelText('playlists')).toBeInTheDocument()
    expect(screen.getByLabelText('top-artists')).toBeInTheDocument()
  })
})
