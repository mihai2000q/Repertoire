import Order from '../../types/Order.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'

const albumsOrders: Order[] = [
  { label: 'Title', property: AlbumProperty.Title, checked: false },
  { label: 'Release Date', property: AlbumProperty.ReleaseDate, nullable: true, checked: false },
  { label: 'Artist', property: AlbumProperty.Artist, nullable: true, checked: false },
  { label: 'Songs', property: AlbumProperty.Songs, checked: false },
  { label: 'Rehearsals', property: AlbumProperty.Rehearsals, checked: false },
  { label: 'Confidence', property: AlbumProperty.Confidence, checked: false },
  { label: 'Progress', property: AlbumProperty.Progress, checked: false },
  { label: 'Last Played', property: AlbumProperty.LastPlayed, nullable: true, checked: false },
  {
    label: 'Creation Date',
    property: AlbumProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: AlbumProperty.LastModified, checked: false }
]

export default albumsOrders
