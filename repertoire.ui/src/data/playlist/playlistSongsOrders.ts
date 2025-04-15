import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const playlistSongsOrders: Order[] = [
  { label: 'Track Number', property: SongProperty.PlaylistTrackNo },
  { label: 'Alphabetically', property: SongProperty.Title },
  {
    label: 'Latest Releases',
    property: SongProperty.ReleaseDate,
    type: OrderType.Descending,
    nullable: true
  },
  {
    label: 'Difficulty',
    property: SongProperty.Difficulty,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.PlaylistTrackNo }]
  },
  {
    label: 'Progress',
    property: SongProperty.Progress,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.PlaylistTrackNo }]
  },
  {
    label: 'Last Played',
    property: SongProperty.LastPlayed,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.PlaylistTrackNo }]
  }
]

export default playlistSongsOrders
