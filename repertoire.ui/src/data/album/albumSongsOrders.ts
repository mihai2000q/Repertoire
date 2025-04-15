import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const albumSongsOrders: Order[] = [
  {
    label: 'Track Number',
    property: SongProperty.AlbumTrackNo,
    type: OrderType.Ascending
  },
  { label: 'Alphabetically', property: SongProperty.Title },
  {
    label: 'Difficulty',
    property: SongProperty.Difficulty,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.AlbumTrackNo }]
  },
  {
    label: 'Rehearsals',
    property: SongProperty.Rehearsals,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.AlbumTrackNo }]
  },
  {
    label: 'Confidence',
    property: SongProperty.Confidence,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.AlbumTrackNo }]
  },
  {
    label: 'Progress',
    property: SongProperty.Progress,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.AlbumTrackNo }]
  },
  {
    label: 'Last Played',
    property: SongProperty.LastPlayed,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.AlbumTrackNo }]
  }
]

export default albumSongsOrders
