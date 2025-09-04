import Playlist from './Playlist.tsx'
import { emptyPlaylist, emptySong, reduxMemoryRouterRender } from '../test-utils.tsx'
import { screen } from '@testing-library/react'
import Song from '../types/models/Song.ts'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { default as PlaylistType } from './../types/models/Playlist.ts'
import { RootState } from '../state/store.ts'
import { MainProvider } from '../context/MainContext.tsx'
import { createRef } from 'react'

describe('Playlist', () => {
  const songs: Song[] = [
    {
      ...emptySong,
      playlistSongId: '1',
      title: 'Song 1'
    },
    {
      ...emptySong,
      playlistSongId: '2',
      title: 'Song 2'
    }
  ]

  const playlist: PlaylistType = {
    ...emptyPlaylist,
    id: '1',
    title: 'Playlist 1'
  }

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
        models: [],
        totalCount: 0
      }
      return HttpResponse.json(response)
    }),
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<PlaylistType> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    }),
    http.get(`/playlists/${playlist.id}`, async () => {
      return HttpResponse.json(playlist)
    }),
    http.get(`/playlists/songs/${playlist.id}`, async () => {
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

  it('should render and display playlist info and songs', async () => {
    const [_, store] = reduxMemoryRouterRender(
      <MainProvider appRef={undefined} scrollRef={createRef()}>
        <Playlist />
      </MainProvider>,
      '/playlist/:id',
      [`/playlist/${playlist.id}`]
    )

    expect(screen.getByTestId('playlist-loader')).toBeInTheDocument()
    expect(await screen.findByLabelText('header-panel-card')).toBeInTheDocument()
    expect(screen.getByLabelText('songs-widget')).toBeInTheDocument()
    expect((store.getState() as RootState).global.documentTitle).toBe(playlist.title)
  })
})
