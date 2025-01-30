import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'

const albumSongsOrders: Order[] = [
  { value: 'album_track_no asc', label: 'Track Number', property: SongProperty.AlbumTrackNo },
  { value: 'title asc', label: 'Alphabetically', property: SongProperty.Title },
  {
    value: 'difficulty desc nulls last, album_track_no asc',
    label: 'Difficulty',
    property: SongProperty.Difficulty
  },
  { value: 'rehearsals desc, album_track_no asc', label: 'Rehearsals', property: SongProperty.Rehearsals },
  { value: 'confidence desc, album_track_no asc', label: 'Confidence', property: SongProperty.Confidence },
  { value: 'progress desc, album_track_no asc', label: 'Progress', property: SongProperty.Progress },
  {
    value: 'last_time_played desc nulls last, album_track_no asc',
    label: 'Last Time Played',
    property: SongProperty.LastTimePlayed
  }
]

export default albumSongsOrders
