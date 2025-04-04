import { emptyArtist, emptySong, reduxRender } from '../../../test-utils.tsx'
import HomeSongCard from './HomeSongCard.tsx'
import Song from '../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'

describe('Home Song Card', () => {
  const song: Song = {
    ...emptySong,
    title: 'Song 1'
  }

  it('should render with minimal info', () => {
    reduxRender(<HomeSongCard song={song} />)

    expect(screen.getByLabelText(`default-icon-${song.title}`)).toBeInTheDocument()
    expect(screen.getByText(song.title)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localSong: Song = {
      ...song,
      imageUrl: 'something.png',
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    reduxRender(<HomeSongCard song={localSong} />)

    expect(screen.getByRole('img', { name: localSong.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localSong.title })).toHaveAttribute(
      'src',
      localSong.imageUrl
    )
    expect(screen.getByText(localSong.title)).toBeInTheDocument()
    expect(screen.getByText(localSong.artist.name)).toBeInTheDocument()
  })

  it('should open song drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRender(<HomeSongCard song={song} />)

    await user.click(screen.getByLabelText(`default-icon-${song.title}`))
    expect((store.getState() as RootState).global.songDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.songDrawer.songId).toBe(song.id)
  })

  it('should open artist drawer by clicking on artist name', async () => {
    const user = userEvent.setup()

    const localSong: Song = {
      ...song,
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    const [_, store] = reduxRender(<HomeSongCard song={localSong} />)

    await user.click(screen.getByText(localSong.artist.name))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localSong.artist.id)
  })
})
