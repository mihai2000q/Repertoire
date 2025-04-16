import Order from '../../types/Order.ts'
import AlbumProperty from '../../utils/enums/AlbumProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

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
  }
]

export default artistAlbumsOrders
