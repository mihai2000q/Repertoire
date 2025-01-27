import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'

const artistSongsOrders: Order[] = [
  {
    value: 'release_date desc nulls last, title asc',
    label: 'Latest Releases',
    property: SongProperty.ReleaseDate
  },
  { value: 'title asc', label: 'Alphabetically', property: SongProperty.Title },
  {
    value: 'difficulty desc nulls last, title asc',
    label: 'Difficulty',
    property: SongProperty.Difficulty
  },
  { value: 'rehearsals desc, title asc', label: 'Rehearsals', property: SongProperty.Rehearsals },
  { value: 'confidence desc, title asc', label: 'Confidence', property: SongProperty.Confidence },
  { value: 'progress desc, title asc', label: 'Progress', property: SongProperty.Progress },
  {
    value: 'last_time_played desc nulls last, title asc',
    label: 'Last Time Played',
    property: SongProperty.LastTimePlayed
  }
]

export default artistSongsOrders
