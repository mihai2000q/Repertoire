import Order from '../../types/Order.ts'
import PlaylistProperty from '../../types/enums/properties/PlaylistProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'

const playlistsOrders: Order[] = [
  { label: 'Title', property: PlaylistProperty.Title, checked: false },
  {
    label: 'Creation Date',
    property: PlaylistProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: PlaylistProperty.LastModified, checked: false },
  { label: 'Songs', property: PlaylistProperty.Songs, checked: false }
]

export default playlistsOrders
