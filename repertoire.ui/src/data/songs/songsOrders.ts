import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const songsOrders: Order[] = [
  { label: 'Title', property: SongProperty.Title, checked: false },
  { label: 'Release Date', property: SongProperty.ReleaseDate, nullable: true, checked: false },
  { label: 'Artist', property: SongProperty.Artist, nullable: true, checked: false },
  { label: 'Album', property: SongProperty.Album, nullable: true, checked: false },
  { label: 'Guitar Tuning', property: SongProperty.GuitarTuning, nullable: true, checked: false },
  { label: 'Difficulty', property: SongProperty.Difficulty, nullable: true, checked: false },
  { label: 'Sections', property: SongProperty.Sections, checked: false },
  { label: 'Solos', property: SongProperty.Solos, checked: false },
  { label: 'Riffs', property: SongProperty.Riffs, checked: false },
  { label: 'Rehearsals', property: SongProperty.Rehearsals, checked: false },
  { label: 'Confidence', property: SongProperty.Confidence, checked: false },
  { label: 'Progress', property: SongProperty.Progress, checked: false },
  { label: 'Last Played', property: SongProperty.LastPlayed, nullable: true, checked: false },
  {
    label: 'Creation Date',
    property: SongProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: SongProperty.LastModified, checked: false }
]

export default songsOrders
