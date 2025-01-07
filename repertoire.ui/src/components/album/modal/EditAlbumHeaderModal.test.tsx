import { reduxRender } from '../../../test-utils.tsx'
import EditAlbumHeaderModal from './EditAlbumHeaderModal.tsx'
import Album from '../../../types/models/Album.ts'

describe('Edit Album Header Modal', () => {
  const album: Album = {
    id: '1',
    title: 'Album 1',
    songs: [],
    createdAt: '',
    updatedAt: '',
    releaseDate: null
  }

  it('should render', () => {
    reduxRender(<EditAlbumHeaderModal opened={true} onClose={() => {}} album={album} />)
  })
})
