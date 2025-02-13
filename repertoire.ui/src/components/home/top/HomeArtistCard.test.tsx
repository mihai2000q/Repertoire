import { emptyArtist, reduxRender } from '../../../test-utils.tsx'
import HomeArtistCard from './HomeArtistCard.tsx'
import Artist from '../../../types/models/Artist.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { RootState } from '../../../state/store.ts'

describe('Home Artist Card', () => {
  const artist: Artist = {
    ...emptyArtist,
    name: 'Artist 1',
    imageUrl: 'something.png'
  }

  it('should render with minimal info', () => {
    reduxRender(<HomeArtistCard artist={artist} />)

    expect(screen.getByRole('img', { name: artist.name })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: artist.name })).toHaveAttribute('src', artist.imageUrl)
    expect(screen.getByText(artist.name)).toBeInTheDocument()
  })

  it('should open artist drawer by clicking on image', async () => {
    const user = userEvent.setup()

    const [_, store] = reduxRender(<HomeArtistCard artist={artist} />)

    await user.click(screen.getByRole('img', { name: artist.name }))
    expect((store.getState() as RootState).global.artistDrawer.open).toBeTruthy()
    expect((store.getState() as RootState).global.artistDrawer.artistId).toBe(artist.id)
  })
})
