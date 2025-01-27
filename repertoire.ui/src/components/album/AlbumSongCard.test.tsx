import { emptySong, reduxRender } from '../../test-utils.tsx'
import AlbumSongCard from './AlbumSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import dayjs from 'dayjs'
import Difficulty from '../../utils/enums/Difficulty.ts'

describe('Album Song Card', () => {
  const song: Song = {
    ...emptySong,
    id: '1',
    title: 'Song 1',
    albumTrackNo: 1
  }

  const emptyOrder: Order = {
    label: '',
    value: ''
  }

  it('should render and display information, when the album is not unknown', () => {
    reduxRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
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

    const [{ rerender }] = reduxRender(
      <AlbumSongCard
        song={localSong}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute('src', localSong.imageUrl)

    const localSongWithAlbum: Song = {
      ...song,
      album: {
        id: '',
        title: '',
        songs: [],
        createdAt: '',
        updatedAt: '',
        imageUrl: 'something-album.png'
      }
    }

    rerender(
      <AlbumSongCard
        song={localSongWithAlbum}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRender(
      <AlbumSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownAlbum={false}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
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

      const [{ rerender }] = reduxRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(
        <AlbumSongCard song={song} handleRemove={() => {}} isUnknownAlbum={false} order={order} />
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

      reduxRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
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

      reduxRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
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

      reduxRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
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

      const [{ rerender }] = reduxRender(
        <AlbumSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownAlbum={false}
          order={order}
        />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('DD MMM YYYY'))
      ).toBeInTheDocument()

      rerender(
        <AlbumSongCard song={song} handleRemove={() => {}} isUnknownAlbum={false} order={order} />
      )

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  describe('on menu', () => {
    it('should display warning modal and remove, when clicking on remove', async () => {
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      reduxRender(
        <AlbumSongCard
          song={song}
          handleRemove={handleRemove}
          isUnknownAlbum={false}
          order={emptyOrder}
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

    reduxRender(
      <AlbumSongCard song={song} handleRemove={() => {}} isUnknownAlbum={true} order={emptyOrder} />
    )

    expect(screen.queryByText(song.albumTrackNo)).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })
})
