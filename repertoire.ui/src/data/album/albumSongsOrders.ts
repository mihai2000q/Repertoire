import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'

const albumSongsOrders: Order[] = [
  { value: 'album_track_no asc', label: 'Track Number', property: SongProperty.AlbumTrackNo },
  { value: 'title asc', label: 'Alphabetically', property: SongProperty.Title },
  { value: 'difficulty desc nulls last', label: 'Difficulty', property: SongProperty.Difficulty },
  { value: 'rehearsals desc', label: 'Rehearsals', property: SongProperty.Rehearsals },
  { value: 'confidence desc', label: 'Confidence', property: SongProperty.Confidence },
  { value: 'progress desc', label: 'Progress', property: SongProperty.Progress },
  { value: 'last_time_played desc nulls last', label: 'Last Time Played', property: SongProperty.LastTimePlayed }
]

export default albumSongsOrders
