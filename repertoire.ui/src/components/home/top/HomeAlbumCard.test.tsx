import { emptyAlbum, emptyArtist, reduxRender } from '../../../test-utils.tsx'
import HomeAlbumCard from './HomeAlbumCard.tsx'
import Album from '../../../types/models/Album.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'

describe('Home Album Card', () => {
  const album: Album = {
    ...emptyAlbum,
    title: 'Album 1'
  }

  it('should render with minimal info', () => {
    reduxRender(<HomeAlbumCard album={album} />)

    expect(screen.getByRole('img', { name: album.title })).toBeInTheDocument()
    expect(screen.getByText(album.title)).toBeInTheDocument()
  })

  it('should render with maximal info', () => {
    const localAlbum: Album = {
      ...album,
      imageUrl: 'something.png',
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    reduxRender(<HomeAlbumCard album={localAlbum} />)

    expect(screen.getByRole('img', { name: localAlbum.title })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: localAlbum.title })).toHaveAttribute(
      'src',
      localAlbum.imageUrl
    )
    expect(screen.getByText(localAlbum.title)).toBeInTheDocument()
    expect(screen.getByText(localAlbum.artist.name)).toBeInTheDocument()
  })

  it('should open album drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRender(<HomeAlbumCard album={album} />)

    await user.click(screen.getByRole('img', { name: album.title }))
    expect((store.getState() as RootState).global.albumDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.albumDrawer.albumId).toBe(album.id)
  })

  it('should open artist drawer by clicking on artist name', async () => {
    const user = userEvent.setup()

    const localAlbum: Album = {
      ...album,
      artist: {
        ...emptyArtist,
        name: 'Artist 1'
      }
    }

    const [_, store] = reduxRender(<HomeAlbumCard album={localAlbum} />)

    await user.click(screen.getByText(localAlbum.artist.name))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(localAlbum.artist.id)
  })
})
