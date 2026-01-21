import Order from '../../types/Order.ts'
import AlbumProperty from '../../types/enums/properties/AlbumProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'

const artistAlbumsOrders: Order[] = [
  {
    label: 'Latest Releases',
    property: AlbumProperty.ReleaseDate,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: AlbumProperty.Title }]
  },
  {
    label: 'Alphabetically',
    property: AlbumProperty.Title,
    type: OrderType.Ascending
  },
  {
    label: 'Rehearsals',
    property: AlbumProperty.Rehearsals,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: AlbumProperty.Title }]
  },
  {
    label: 'Confidence',
    property: AlbumProperty.Confidence,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: AlbumProperty.Title }]
  },
  {
    label: 'Progress',
    property: AlbumProperty.Progress,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: AlbumProperty.Title }]
  },
  {
    label: 'Last Played',
    property: AlbumProperty.LastPlayed,
    type: OrderType.Descending,
    nullable: true,
    thenBy: [{ property: AlbumProperty.Title }]
  }
]

export default artistAlbumsOrders
