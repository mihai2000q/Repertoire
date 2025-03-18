import { emptyOrder, emptySong, reduxRouterRender } from '../../test-utils.tsx'
import AlbumSongCard from './AlbumSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import SongProperty from '../../utils/enums/SongProperty.ts'
import dayjs from 'dayjs'
import Difficulty from '../../utils/enums/Difficulty.ts'
import { expect } from 'vitest'

describe('Album Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1',
    albumTrackNo: 1
  }

  it('should render and display information, when the album is not unknown', () => {
    reduxRouterRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )

    expect(screen.getByText(song.albumTrackNo)).toBeInTheDocument()
    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
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
        handleRemove={() => {}}
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
        handleRemove={() => {}}
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
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(
        <AlbumSongCard
          song={song}
          handleRemove={() => {}}
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
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getAllByText(localSong.rehearsals)).toHaveLength(2) // the one visible and the one in the tooltip
      expect(screen.getAllByText(localSong.rehearsals)[0]).toBeVisible()
      expect(screen.getAllByText(localSong.rehearsals)[1]).not.toBeVisible()
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
          handleRemove={() => {}}
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
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the date, when it is last time played', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.LastTimePlayed
      }

      const localSong = {
        ...song,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
          isDragging={false}
        />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('D MMM YYYY'))
      ).toBeInTheDocument()

      rerender(
        <AlbumSongCard
          song={song}
          handleRemove={() => {}}
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

    reduxRouterRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-card-${song.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
        isDragging={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /partial rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /perfect rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <AlbumSongCard
          song={song}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/song/${song.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and remove, when clicking on remove', async () => {
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      reduxRouterRender(
        <AlbumSongCard
          song={song}
          handleRemove={handleRemove}
          isUnknownAlbum={false}
          order={emptyOrder}
          isDragging={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /remove/i }))

      expect(await screen.findByRole('dialog', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /remove song/i })).toBeInTheDocument()
      expect(screen.getByText(/are you sure/i)).toBeInTheDocument()

      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(handleRemove).toHaveBeenCalledOnce()
    })
  })

  it('should not display the tracking number and some menu options, when the album is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={true}
        order={emptyOrder}
        isDragging={false}
      />
    )

    expect(screen.queryByText(song.albumTrackNo)).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })
})
