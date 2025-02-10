import { emptyArtist, mantineRender } from '../../../test-utils.tsx'
import ArtistInfoModal from './ArtistInfoModal.tsx'
import Artist from '../../../types/models/Artist.ts'
import { screen } from '@testing-library/react'
import dayjs from 'dayjs'

describe('Artist Info Modal', () => {
  const artist: Artist = {
    ...emptyArtist,
    id: '1',
    name: 'Artist 1',
    createdAt: '2024-11-25T22:00:00',
    updatedAt: '2024-12-12T05:00:00'
  }

  it('should render', () => {
    mantineRender(<ArtistInfoModal opened={true} onClose={() => {}} artist={artist} />)

    expect(screen.getByRole('dialog', { name: /artist info/i })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /artist info/i })).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(artist.createdAt).format('DD MMMM YYYY, HH:mm'))
    ).toBeInTheDocument()
    expect(
      screen.getByText(dayjs(artist.updatedAt).format('DD MMMM YYYY, HH:mm'))
    ).toBeInTheDocument()
  })
})
