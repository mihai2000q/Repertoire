import Order from '../../types/Order.ts'
import PlaylistProperty from '../../utils/enums/PlaylistProperty.ts'
import OrderType from '../../utils/enums/OrderType.ts'

const playlistsOrders: Order[] = [
  { label: 'Title', property: PlaylistProperty.Title, checked: false },
  { label: 'Songs', property: PlaylistProperty.Songs, checked: false },
  {
    label: 'Creation Date',
    property: PlaylistProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: PlaylistProperty.LastModified, checked: false }
]

export default playlistsOrders
