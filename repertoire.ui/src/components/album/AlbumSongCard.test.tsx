import { emptyOrder, emptySong, reduxRouterRender, withToastify } from '../../test-utils.tsx'
import AlbumSongCard from './AlbumSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { fireEvent, screen, within } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import SongProperty from '../../types/enums/SongProperty.ts'
import dayjs from 'dayjs'
import Difficulty from '../../types/enums/Difficulty.ts'
import { expect } from 'vitest'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { RemoveSongsFromAlbumRequest } from '../../types/requests/AlbumRequests.ts'
import WithTotalCountResponse from '../../types/responses/WithTotalCountResponse.ts'
import Playlist from '../../types/models/Playlist.ts'

describe('Album Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1',
    albumTrackNo: 1
  }

  const handlers = [
    http.get('/playlists', async () => {
      const response: WithTotalCountResponse<Playlist> = { models: [], totalCount: 0 }
      return HttpResponse.json(response)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    server.resetHandlers()
    window.location.pathname = '/'
  })

  afterAll(() => server.close())

  it('should render and display information, when the album is not unknown', () => {
    reduxRouterRender(
      <AlbumSongCard
        song={song}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )

    expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(
      <AlbumSongCard
        song={localSong}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', localSong.imageUrl)

    const albumImageUrl = 'something.png'

    rerender(
      <AlbumSongCard
        song={song}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
        albumImageUrl={albumImageUrl}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', albumImageUrl)
  })

  describe('on order property change', () => {
    it('should display the difficulty bar, when it is difficulty', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Difficulty
      }

      const localSong = {
        ...song,
        difficulty: Difficulty.Easy
      }

      const [{ rerender }] = reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(
        <AlbumSongCard
          song={song}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()
    })

    it('should display the rehearsals, when it is rehearsals', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Rehearsals
      }

      const localSong = {
        ...song,
        rehearsals: 34
      }

      reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getAllByText(localSong.rehearsals)).toHaveLength(2) // the one visible and the one in the tooltip
    })

    it('should display the confidence bar, when it is confidence', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Confidence
      }

      const localSong = {
        ...song,
        confidence: 23
      }

      reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    })

    it('should display the progress bar, when it is progress', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.Progress
      }

      const localSong = {
        ...song,
        progress: 123
      }

      reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the date, when it is last played', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.LastPlayed
      }

      const localSong = {
        ...song,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
      ).toBeInTheDocument()

      rerender(
        <AlbumSongCard
          song={song}
          albumId={''}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <AlbumSongCard
        song={song}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-card-${song.title}`)
    })
    expectDefaultMenu()

    rerender(
      <AlbumSongCard
        song={{
          ...song,
          songsterrLink: 'some-link',
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI'
        }}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await expectFullMenu()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = reduxRouterRender(
      <AlbumSongCard
        song={song}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await user.click(screen.getByRole('button', { name: 'more-menu' }))
    expectDefaultMenu()

    rerender(
      <AlbumSongCard
        song={{
          ...song,
          songsterrLink: 'some-link',
          youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI'
        }}
        albumId={''}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )
    await expectFullMenu()
  })

  function expectDefaultMenu() {
    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add to playlist/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove from album/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /open links/i })).toBeDisabled()
  }

  async function expectFullMenu() {
    const user = userEvent.setup()

    expect(screen.getByRole('menuitem', { name: /open links/i })).not.toBeDisabled()
    await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
    expect(screen.getByRole('menuitem', { name: /songsterr/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).toBeInTheDocument()

    expect(screen.getByRole('menuitem', { name: /songsterr/i })).not.toBeDisabled()
    expect(screen.getByRole('menuitem', { name: /youtube/i })).not.toBeDisabled()
  }

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumSongCard
          song={song}
          albumId={''}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/song/${song.id}`)
    })

    it('should open youtube when clicking on open links => open youtube', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, youtubeLink: 'https://www.youtube.com/watch?v=tAGnKpE4NCI' }

      server.use(
        http.get(
          localSong.youtubeLink
            .replace('youtube', 'youtube-nocookie')
            .replace('watch?v=', 'embed/'),
          () => {
            return HttpResponse.json({ message: 'it worked' })
          }
        )
      )

      reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
      fireEvent.click(screen.getByRole('menuitem', { name: /youtube/i }))

      expect(await screen.findByRole('dialog', { name: song.title })).toBeInTheDocument()
    })

    it('should be able to open songsterr in browser when clicking on open links => open songsterr', async () => {
      const user = userEvent.setup()

      const localSong = { ...song, songsterrLink: 'some-link' }

      reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          albumId={''}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.hover(screen.getByRole('menuitem', { name: /open links/i }))
      expect(
        within(screen.getByRole('link', { name: /songsterr/i })).getByRole('menuitem', {
          name: /songsterr/i
        })
      ).toBeInTheDocument()
      expect(screen.getByRole('link', { name: /songsterr/i })).toBeExternalLink(
        localSong.songsterrLink
      )
    })

    it('should display warning modal and remove song from album, when clicking on remove song', async () => {
      const user = userEvent.setup()

      let capturedRequest: RemoveSongsFromAlbumRequest
      server.use(
        http.put('/albums/remove-songs', async (req) => {
          capturedRequest = (await req.request.json()) as RemoveSongsFromAlbumRequest
          return HttpResponse.json()
        })
      )

      const albumId = 'some-album-id'

      reduxRouterRender(
        withToastify(
          <AlbumSongCard
            song={song}
            albumId={albumId}
            isUnknownAlbum={false}
            order={emptyOrder}
            isDragging={false}
          />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove from album/i }))

      expect(await screen.findByRole('dialog', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(capturedRequest).toStrictEqual({
        id: albumId,
        songIds: [song.id]
      })
      expect(await screen.findByText(new RegExp(`${song.title} removed`, 'i'))).toBeInTheDocument()
    })

    it('should display warning modal and delete album, when clicking on delete', async () => {
      const user = userEvent.setup()

      server.use(
        http.delete(`/songs/${song.id}`, () => {
          return HttpResponse.json()
        })
      )

      reduxRouterRender(
        <AlbumSongCard
          song={song}
          albumId={''}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))
    })
  })

  it('should not display the tracking number and some menu options, when the album is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongCard
        song={song}
        albumId={''}
        isUnknownAlbum={true}
        order={emptyOrder}
        isDragging={false}
      />
    )

    expect(screen.queryByText(song.albumTrackNo)).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove from album/i })).not.toBeInTheDocument()
  })
})
