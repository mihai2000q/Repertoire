import { mantineRender } from '../../../test-utils.tsx'
import AlbumInfoModal from './AlbumInfoModal.tsx'
import Album from '../../../types/models/Album.ts'
import {screen} from "@testing-library/react";

describe('Album Info Modal', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    songs: [],
    createdAt: '2024-11-25T22:00:00',
    updatedAt: '2024-12-12T05:00:00'
  }

  it('should render', () => {
    mantineRender(<AlbumInfoModal opened={true} onClose={() => {}} album={album} />)

    expect(screen.getByRole('heading', { name: /album info/i })).toBeInTheDocument()
    expect(screen.getByText('25 November 2024, 22:00')).toBeInTheDocument()
    expect(screen.getByText('12 December 2024, 05:00')).toBeInTheDocument()
  })
})
