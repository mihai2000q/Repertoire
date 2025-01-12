import { reduxRender } from '../../test-utils.tsx'
import PlaylistSongsCard from './PlaylistSongsCard.tsx'
import Song from '../../types/models/Song.ts'
import Playlist from '../../types/models/Playlist.ts'
import { screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import { RemoveSongsFromPlaylistRequest } from "../../types/requests/PlaylistRequests.ts";

describe('Playlist Songs Card', () => {
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

  const playlist: Playlist = {
    id: '1',
    title: 'Song 1',
    description: '',
    createdAt: '',
    updatedAt: '',
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

  const handlers = [
    http.get('/songs', async () => {
      const response: WithTotalCountResponse<Song> = {
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

  it("should render and display playlist's songs", () => {
    reduxRender(<PlaylistSongsCard playlist={playlist} />)

    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(playlist.songs.length)
    playlist.songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()
  })

  it('should display menu', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistSongsCard playlist={playlist} />)

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add songs/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open add playlist songs modal', async () => {
      const user = userEvent.setup()

      // Act
      reduxRender(<PlaylistSongsCard playlist={playlist} />)

      // Assert
      await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /add songs/i }))

      expect(screen.getByRole('heading', { name: /add playlist songs/i })).toBeInTheDocument()
    })
  })

  it('should display new song card when there are no playlist songs and open Add playlist songs modal', async () => {
    const user = userEvent.setup()

    reduxRender(<PlaylistSongsCard playlist={{ ...playlist, songs: [] }} />)

    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(screen.getByRole('heading', { name: /add playlist songs/i })).toBeInTheDocument()
  })

  it('should send \'remove songs from playlist request\' when clicking on the more menu of a song card', async () => {
    const user = userEvent.setup()

    const song = playlist.songs[0]

    let capturedRequest: RemoveSongsFromPlaylistRequest
    server.use(
      http.put('/playlists/remove-songs', async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromPlaylistRequest
        return HttpResponse.json()
      })
    )

    reduxRender(<PlaylistSongsCard playlist={playlist} />)

    const songCard1 = screen.getByLabelText(`song-card-${song.title}`)

    await user.click(within(songCard1).getByRole('button', { name: 'more-menu' }))
    await user.click(screen.getByRole('menuitem', { name: /remove/i }))
    await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

    expect(capturedRequest.id).toBe(playlist.id)
    expect(capturedRequest.songIds).toStrictEqual([song.id])
  })
})
