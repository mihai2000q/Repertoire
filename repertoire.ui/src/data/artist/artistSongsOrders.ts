import Order from '../../types/Order.ts'
import SongProperty from '../../utils/enums/SongProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const artistSongsOrders: Order[] = [
  {
    label: 'Latest Releases',
    property: SongProperty.ReleaseDate,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  },
  {
    label: 'Alphabetically',
    property: SongProperty.Title,
    type: OrderType.Ascending
  },
  {
    label: 'Difficulty',
    property: SongProperty.Difficulty,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  },
  {
    label: 'Rehearsals',
    property: SongProperty.Rehearsals,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  },
  {
    label: 'Confidence',
    property: SongProperty.Confidence,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  },
  {
    label: 'Progress',
    property: SongProperty.Progress,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  },
  {
    label: 'Last Time Played',
    property: SongProperty.LastTimePlayed,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  }
]

export default artistSongsOrders
