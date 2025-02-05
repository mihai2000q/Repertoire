import Order from '../../types/Order.ts'
import AlbumProperty from '../../utils/enums/AlbumProperty.ts'

const artistAlbumsOrders: Order[] = [
  {
    value: 'release_date desc nulls last, title asc',
    label: 'Latest Releases',
    property: AlbumProperty.ReleaseDate
  },
  { value: 'title asc', label: 'Alphabetically', property: AlbumProperty.Title }
]

export default artistAlbumsOrders
