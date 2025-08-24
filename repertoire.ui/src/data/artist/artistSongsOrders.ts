import Order from '../../types/Order.ts'
import SongProperty from '../../types/enums/properties/SongProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'

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
    label: 'Last Played',
    property: SongProperty.LastPlayed,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: SongProperty.Title }]
  }
]

export default artistSongsOrders
