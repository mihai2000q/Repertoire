import Order from '../../types/Order.ts'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import OrderType from '../../types/enums/OrderType.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'

const artistsOrders: Order[] = [
  { label: 'Name', property: ArtistProperty.Name, checked: false },
  {
    label: 'Creation Date',
    property: ArtistProperty.CreationDate,
    type: OrderType.Descending,
    checked: true
  },
  { label: 'Last Modified', property: ArtistProperty.LastModified, checked: false },
  { label: 'Band', property: ArtistProperty.Band, checked: false },
  { label: 'Band Members', property: ArtistProperty.BandMembers, checked: false },
  { label: 'Albums', property: ArtistProperty.Albums, checked: false },
  { label: 'Songs', property: ArtistProperty.Songs, checked: false },
  { label: 'Rehearsals', property: ArtistProperty.Rehearsals, checked: false },
  { label: 'Confidence', property: ArtistProperty.Confidence, checked: false },
  { label: 'Progress', property: ArtistProperty.Progress, checked: false },
  { label: 'Last Played', property: AlbumProperty.LastPlayed, nullable: true, checked: false }
]

export default artistsOrders
