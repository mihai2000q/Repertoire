import { emptyOrder, reduxRouterRender } from '../../test-utils.tsx'
import ArtistSongCard from './ArtistSongCard.tsx'
import Song from '../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import Album from 'src/types/models/Album.ts'
import { RootState } from '../../state/store.ts'
import dayjs from 'dayjs'
import SongProperty from '../../utils/enums/SongProperty.ts'
import Difficulty from '../../utils/enums/Difficulty.ts'
import { expect } from 'vitest'

describe('Artist Song Card', () => {
  const song: Song = {
    id: '1',
    title: 'Song 1',
    description: '',
    isRecorded: false,
    rehearsals: 0,
    confidence: 0,
    progress: 0,
    sections: [],
    createdAt: '',
    updatedAt: ''
  }

  const album: Album = {
    id: '1',
    title: 'Album 1',
    createdAt: '',
    updatedAt: '',
    songs: []
  }

  it('should render and display minimal information', async () => {
    reduxRouterRender(
      <ArtistSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal information', async () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png',
      album: album
    }

    reduxRouterRender(
      <ArtistSongCard
        song={localSong}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.album.title)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it("should display song's image if the song has one, if not the album's image", () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png'
    }

    const [{ rerender }] = reduxRouterRender(
      <ArtistSongCard
        song={localSong}
        handleRemove={() => {}}
        isUnknownArtist={false}
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
      <ArtistSongCard
        song={localSongWithAlbum}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    expect(screen.getByRole('img', { name: song.title })).toHaveAttribute(
      'src',
      localSongWithAlbum.album.imageUrl
    )
  })

  describe('on order property change', () => {
    it('should display the release date, when it is release date', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.ReleaseDate
      }

      const localSong = {
        ...song,
        releaseDate: '2024-10-12T10:30'
      }

      reduxRouterRender(
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={order}
        />
      )

      expect(
        screen.getByText(dayjs(localSong.releaseDate).format('D MMM YYYY'))
      ).toBeInTheDocument()
    })

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
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={order}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'difficulty' })).toBeInTheDocument()

      rerender(
        <ArtistSongCard song={song} handleRemove={() => {}} isUnknownArtist={false} order={order} />
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
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
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

      reduxRouterRender(
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
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

      reduxRouterRender(
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={order}
        />
      )

      expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    })

    it('should display the last time played date, when it is last time played', () => {
      const order = {
        ...emptyOrder,
        property: SongProperty.LastTimePlayed
      }

      const localSong = {
        ...song,
        lastTimePlayed: '2024-10-12T10:30'
      }

      const [{ rerender }] = reduxRouterRender(
        <ArtistSongCard
          song={localSong}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={order}
        />
      )

      expect(
        screen.getByText(dayjs(localSong.lastTimePlayed).format('D MMM YYYY'))
      ).toBeInTheDocument()

      rerender(
        <ArtistSongCard song={song} handleRemove={() => {}} isUnknownArtist={false} order={order} />
      )

      expect(screen.getByText(/never/i)).toBeInTheDocument()
    })
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-card-${song.title}`)
    })

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /remove/i })).toBeInTheDocument()
  })

  it('should display less information on the menu when the artist is unknown', async () => {
    const user = userEvent.setup()

    reduxRouterRender(
      <ArtistSongCard
        song={song}
        handleRemove={() => {}}
        isUnknownArtist={true}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /view details/i })).toBeInTheDocument()
    expect(screen.queryByRole('menuitem', { name: /remove/i })).not.toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should navigate to song when clicking on view details', async () => {
      const user = userEvent.setup()

      reduxRouterRender(
        <ArtistSongCard
          song={song}
          handleRemove={() => {}}
          isUnknownArtist={false}
          order={emptyOrder}
        />
      )

      await user.click(await screen.findByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /view details/i }))

      expect(window.location.pathname).toBe(`/song/${song.id}`)

      // restore
      window.location.pathname = '/'
    })

    it('should display warning modal and remove song, when clicking on remove', async () => {
      const user = userEvent.setup()

      const handleRemove = vitest.fn()

      reduxRouterRender(
        <ArtistSongCard
          song={song}
          handleRemove={handleRemove}
          isUnknownArtist={false}
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

  it('should open album drawer on album title click', async () => {
    const user = userEvent.setup()

    const localSong = {
      ...song,
      album: album
    }

    const [_, store] = reduxRouterRender(
      <ArtistSongCard
        song={localSong}
        handleRemove={() => {}}
        isUnknownArtist={false}
        order={emptyOrder}
      />
    )

    await user.click(screen.getByText(localSong.album.title))

    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(localSong.album.id)
  })
})
