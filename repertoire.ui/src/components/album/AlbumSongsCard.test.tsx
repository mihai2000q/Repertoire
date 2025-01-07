import { reduxRender } from '../../test-utils.tsx'
import AlbumSongsCard from './AlbumSongsCard.tsx'
import Song from '../../types/models/Song.ts'
import Album from '../../types/models/Album.ts'
import {screen, within} from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { http, HttpResponse } from 'msw'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import { setupServer } from 'msw/node'
import {RemoveSongsFromAlbumRequest} from "../../types/requests/AlbumRequests.ts";

describe('Album Songs Card', () => {
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

  const album: Album = {
    id: '1',
    title: 'Song 1',
    createdAt: '',
    updatedAt: '',
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

  it("should render and display album's songs when the album is not unknown", async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={[]} isUnknownAlbum={false} />
    )

    // Assert
    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(album.songs.length)
    album.songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.queryByLabelText('new-song-card')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.getByRole('menuitem', { name: /add existing songs/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should render and display the songs when the album is unknown', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={songs} isUnknownAlbum={true} />
    )

    // Assert
    expect(screen.getByText(/songs/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'songs-more-menu' })).toBeInTheDocument()
    expect(screen.getAllByLabelText(/song-card-/)).toHaveLength(songs.length)
    songs.forEach((song) =>
      expect(screen.getByLabelText(`song-card-${song.title}`)).toBeInTheDocument()
    )
    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /add existing songs/i })).not.toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should open add existing songs modal when clicking on the menu item', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={[]} isUnknownAlbum={false} />
    )

    // Assert
    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
    await user.click(screen.getByRole('menuitem', { name: /add existing songs/i }))

    expect(screen.getByRole('heading', { name: /add existing songs/i })).toBeInTheDocument()
  })

  it('should open add existing songs modal when clicking on the menu item', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={[]} isUnknownAlbum={false} />
    )

    // Assert
    await user.click(screen.getByRole('button', { name: 'songs-more-menu' }))
    await user.click(screen.getByRole('menuitem', { name: /add new song/i }))

    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display new song card when there are no album songs and open Add existing songs modal', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={{ ...album, songs: [] }} songs={[]} isUnknownAlbum={false} />
    )

    // Assert
    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(screen.getByRole('heading', { name: /add existing songs/i })).toBeInTheDocument()
  })

  it('should display new song card when the album is unknown and open Add new song modal', async () => {
    // Arrange
    const user = userEvent.setup()

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={songs} isUnknownAlbum={true} />
    )

    // Assert
    expect(screen.getByLabelText('new-song-card')).toBeInTheDocument()

    await user.click(screen.getByLabelText('new-song-card'))

    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()
  })

  it('should display send remove song from album request when clicking on the more menu of a song card', async () => {
    // Arrange
    const user = userEvent.setup()

    const song = album.songs[0]

    let capturedRequest: RemoveSongsFromAlbumRequest
    server.use(
      http.put('/albums/remove-songs', async (req) => {
        capturedRequest = (await req.request.json()) as RemoveSongsFromAlbumRequest
        return HttpResponse.json()
      })
    )

    // Act
    reduxRender(
      <AlbumSongsCard album={album} songs={songs} isUnknownAlbum={false} />
    )

    // Assert
    const songCard1 = screen.getByLabelText(`song-card-${song.title}`)

    await user.click(within(songCard1).getByRole('button', { name: 'more-menu' }))
    await user.click(screen.getByRole('menuitem', { name: /remove/i }))
    await user.click(screen.getByRole('button', { name: /yes/i })) // warning modal

    expect(capturedRequest.id).toBe(album.id)
    expect(capturedRequest.songIds).toStrictEqual([song.id])
  })
})
