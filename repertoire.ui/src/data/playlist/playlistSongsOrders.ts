import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'

const playlistSongsOrders: Order[] = [
  { value: 'song_track_no asc', label: 'Track Number', property: SongProperty.PlaylistTrackNo },
  { value: 'release_date desc', label: 'Latest Releases', property: SongProperty.ReleaseDate },
  { value: 'title asc', label: 'Alphabetically' },
  { value: 'release_date asc', label: 'First Releases', property: SongProperty.ReleaseDate }
]

export default playlistSongsOrders
